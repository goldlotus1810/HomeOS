#!/usr/bin/env python3
"""
SILK — Sổ Cái Tự Động cho WORLDTREE
=====================================
Tự động đếm, ghi lại mọi thay đổi trong WORLDTREE.
Mỗi khi thêm file/folder → tự cập nhật MANIFEST.md

Usage:
    python3 silk.py                    # scan once, update MANIFEST
    python3 silk.py --watch            # watch mode (auto-update on change)
    python3 silk.py --watch --interval 5   # check every 5 seconds
    python3 silk.py --diff             # show what changed since last scan
    python3 silk.py --status           # quick summary only
"""

import os, sys, json, time, hashlib, datetime, argparse
from pathlib import Path
from collections import defaultdict

# ── Paths ──────────────────────────────────────────────────────
SILK_DIR    = Path(__file__).parent                   # /tmp/SILK/
WORLDTREE   = SILK_DIR.parent / "WORLDTREE"           # /tmp/WORLDTREE/
MANIFEST    = WORLDTREE / "MANIFEST.md"               # cuốn sổ cái
LEDGER_JSON = SILK_DIR / "ledger.json"                # snapshot nội bộ
CHANGELOG   = SILK_DIR / "CHANGELOG.md"               # lịch sử thay đổi

UNICODE_17_TOTAL  = 159_801
UNICODE_17_BLOCKS = 338

# ── Scan ──────────────────────────────────────────────────────
def scan_worldtree():
    """Đọc toàn bộ WORLDTREE, trả về snapshot dict."""
    if not WORLDTREE.exists():
        raise FileNotFoundError(f"WORLDTREE not found at {WORLDTREE}")

    snap = {
        "timestamp": datetime.datetime.utcnow().isoformat() + "Z",
        "dirs": {},    # path → child_count
        "files": {},   # path → size
        "blocks": {},  # block_name → {files, chars, folder}
    }

    total_dirs = 0
    total_files = 0
    total_bytes = 0

    for root, dirs, files in os.walk(WORLDTREE):
        dirs.sort()
        rel_root = Path(root).relative_to(WORLDTREE)
        snap["dirs"][str(rel_root)] = len(dirs) + len(files)
        total_dirs += 1

        for fname in sorted(files):
            fpath = Path(root) / fname
            rel   = str(fpath.relative_to(WORLDTREE))
            size  = fpath.stat().st_size
            snap["files"][rel] = size
            total_files += 1
            total_bytes += size

    snap["total_dirs"]  = total_dirs
    snap["total_files"] = total_files
    snap["total_bytes"] = total_bytes

    # ── Block stats từ INDEX.md ──
    for ipath in sorted(WORLDTREE.rglob("INDEX.md")):
        folder = str(ipath.parent.relative_to(WORLDTREE))
        try:
            text = ipath.read_text(encoding="utf-8")
            # parse "Assigned characters: N"
            assigned = 0
            for line in text.splitlines():
                if line.startswith("Assigned characters:"):
                    assigned = int(line.split(":")[1].split("|")[0].strip().replace(",",""))
                    break
            # block name = first # heading
            bname = ""
            for line in text.splitlines():
                if line.startswith("# "):
                    bname = line[2:].strip()
                    break
            if bname:
                snap["blocks"][bname] = {
                    "chars": assigned,
                    "folder": folder,
                }
        except Exception:
            pass

    return snap


# ── Diff ──────────────────────────────────────────────────────
def diff_snapshots(old, new):
    """So sánh 2 snapshots, trả về dict các thay đổi."""
    old_files = set(old.get("files", {}).keys())
    new_files = set(new.get("files", {}).keys())
    old_dirs  = set(old.get("dirs",  {}).keys())
    new_dirs  = set(new.get("dirs",  {}).keys())

    return {
        "added_files":   sorted(new_files - old_files),
        "removed_files": sorted(old_files - new_files),
        "added_dirs":    sorted(new_dirs  - old_dirs),
        "removed_dirs":  sorted(old_dirs  - new_dirs),
        "file_delta":    new.get("total_files", 0) - old.get("total_files", 0),
        "dir_delta":     new.get("total_dirs",  0) - old.get("total_dirs",  0),
        "bytes_delta":   new.get("total_bytes", 0) - old.get("total_bytes", 0),
    }


# ── Build MANIFEST.md ─────────────────────────────────────────
def build_manifest(snap, changes=None):
    """Render MANIFEST.md từ snapshot hiện tại."""
    ts      = snap["timestamp"]
    n_files = snap["total_files"]
    n_dirs  = snap["total_dirs"]
    n_bytes = snap["total_bytes"]
    n_chars = sum(b["chars"] for b in snap["blocks"].values())
    missing = UNICODE_17_TOTAL - n_chars
    cov     = n_chars / UNICODE_17_TOTAL * 100

    mb = n_bytes / 1024 / 1024

    now_str = datetime.datetime.utcnow().strftime("%Y-%m-%d %H:%M UTC")

    # ── Change summary ──
    change_section = ""
    if changes:
        af = changes["added_files"]
        rf = changes["removed_files"]
        ad = changes["added_dirs"]
        rd = changes["removed_dirs"]
        fd = changes["file_delta"]
        dd = changes["dir_delta"]
        bd = changes["bytes_delta"]

        lines = [
            "## LAST CHANGE",
            "",
            f"    Recorded:      {now_str}",
            f"    Files Δ:       {fd:+d}  ({len(af)} added, {len(rf)} removed)",
            f"    Folders Δ:     {dd:+d}  ({len(ad)} added, {len(rd)} removed)",
            f"    Size Δ:        {bd:+,} bytes",
        ]
        if ad:
            lines.append("")
            lines.append("    New folders:")
            for d in ad[:20]:
                lines.append(f"      + {d}")
            if len(ad) > 20:
                lines.append(f"      ... and {len(ad)-20} more")
        if af:
            lines.append("")
            lines.append("    New files (first 30):")
            for f in af[:30]:
                lines.append(f"      + {f}")
            if len(af) > 30:
                lines.append(f"      ... and {len(af)-30} more")
        change_section = "\n".join(lines) + "\n\n"

    # ── Top blocks ──
    top_blocks = sorted(snap["blocks"].items(), key=lambda x: x[1]["chars"], reverse=True)

    block_rows = []
    for bname, bdata in top_blocks[:60]:
        chars  = bdata["chars"]
        folder = bdata["folder"]
        block_rows.append(f"| {bname:<45} | {chars:>7,} | {folder} |")

    block_table = "\n".join(block_rows)

    # ── Dir tree (2 levels) ──
    tree_lines = []
    top_level = sorted([d for d in snap["dirs"].keys()
                        if d != "." and "/" not in d and d != ""])
    for tl in top_level:
        child_count = snap["dirs"].get(tl, 0)
        tree_lines.append(f"    {tl}/  ({child_count} items)")
        subs = sorted([d for d in snap["dirs"].keys()
                       if d.startswith(tl + "/") and d.count("/") == 1])
        for s in subs[:15]:
            sc = snap["dirs"].get(s, 0)
            tree_lines.append(f"      {s.split('/')[-1]}/  ({sc} items)")
        if len(subs) > 15:
            tree_lines.append(f"      ... and {len(subs)-15} more subdirectories")
    tree_str = "\n".join(tree_lines)

    manifest = f"""# WORLDTREE — MANIFEST
> *Sổ cái tự động do SILK quản lý. Cập nhật mỗi khi có thay đổi trong WORLDTREE.*
> *Không chỉnh sửa tay file này.*

---

## SNAPSHOT

    Last updated:  {now_str}
    Scan time:     {ts}

    Total files:   {n_files:,}
    Total folders: {n_dirs:,}
    Total size:    {mb:.1f} MB  ({n_bytes:,} bytes)

## UNICODE COVERAGE

    Standard:      Unicode 17.0.0  (released 2025-09-09)
    Total 17.0:    {UNICODE_17_TOTAL:,} characters
    In WORLDTREE:  {n_chars:,} characters
    Missing:       {missing:,} characters
    Coverage:      {cov:.1f}%

    Missing mainly:
      - CJK Ext J (323B0–3347F): ~4,298  new in 17.0, needs UnicodeData.txt
      - CJK Ext B–H:             ~18,000 ideographs beyond Python 3.12 data
      - New 17.0 scripts:        Tai Yo, Beria Erfe (partial coverage)

    To reach 100%:
      wget https://unicode.org/Public/17.0.0/ucd/UnicodeData.txt
      python3 SILK/build_missing.py UnicodeData.txt

{change_section}## WORLDTREE STRUCTURE

{tree_str}

## TOP BLOCKS BY CHARACTER COUNT

| Block                                             |   Chars | Folder |
|---------------------------------------------------|---------|--------|
{block_table}

## RULES

    WORLDTREE is append-only:
      ✓  Add files and folders freely
      ✗  Never rename existing files or folders
      ✗  Never delete files or folders
      ✗  Never edit existing character (.md) files

    SILK (this ledger) auto-updates MANIFEST.md on every change.
    SILK records all changes in SILK/CHANGELOG.md.

## SILK LEDGER FILES

    SILK/silk.py        — the ledger engine (run to update)
    SILK/ledger.json    — current snapshot (machine-readable)
    SILK/CHANGELOG.md   — full history of all changes
    WORLDTREE/MANIFEST.md — this file (human-readable summary)
"""
    return manifest


# ── CHANGELOG append ──────────────────────────────────────────
def append_changelog(changes, snap):
    now_str = datetime.datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S UTC")
    af = changes["added_files"]
    rf = changes["removed_files"]
    ad = changes["added_dirs"]
    rd = changes["removed_dirs"]
    fd = changes["file_delta"]

    if fd == 0 and changes["dir_delta"] == 0:
        return  # no change, skip

    entry = [
        f"\n## {now_str}",
        f"",
        f"Files: {fd:+d} | Dirs: {changes['dir_delta']:+d} | Size: {changes['bytes_delta']:+,} bytes",
        "",
    ]
    if ad:
        entry.append("Added folders:")
        for d in ad: entry.append(f"  + {d}")
    if rd:
        entry.append("Removed folders:")
        for d in rd: entry.append(f"  - {d}")
    if af:
        entry.append(f"Added files ({len(af)}):")
        for f in af[:50]: entry.append(f"  + {f}")
        if len(af) > 50: entry.append(f"  ... and {len(af)-50} more")
    if rf:
        entry.append(f"Removed files ({len(rf)}):")
        for f in rf[:20]: entry.append(f"  - {f}")

    entry.append("")
    entry_text = "\n".join(entry)

    # Append hoặc tạo mới
    if CHANGELOG.exists():
        with open(CHANGELOG, "a", encoding="utf-8") as f:
            f.write(entry_text)
    else:
        header = "# SILK CHANGELOG\n\nLịch sử mọi thay đổi trong WORLDTREE.\n"
        CHANGELOG.write_text(header + entry_text, encoding="utf-8")


# ── Save / Load ledger ────────────────────────────────────────
def load_ledger():
    if LEDGER_JSON.exists():
        return json.loads(LEDGER_JSON.read_text(encoding="utf-8"))
    return {}

def save_ledger(snap):
    LEDGER_JSON.write_text(
        json.dumps(snap, indent=2, ensure_ascii=False),
        encoding="utf-8"
    )


# ── Main update cycle ─────────────────────────────────────────
def update(verbose=True):
    old_snap = load_ledger()
    new_snap = scan_worldtree()
    changes  = diff_snapshots(old_snap, new_snap) if old_snap else None

    # Render MANIFEST
    manifest_text = build_manifest(new_snap, changes)
    MANIFEST.write_text(manifest_text, encoding="utf-8")

    # Log changes
    if changes:
        append_changelog(changes, new_snap)

    # Save new snapshot
    save_ledger(new_snap)

    if verbose:
        ts = datetime.datetime.utcnow().strftime("%H:%M:%S")
        n  = new_snap["total_files"]
        nc = sum(b["chars"] for b in new_snap["blocks"].values())
        print(f"[{ts}] SILK ✓  {n:,} files | {nc:,} chars", end="")
        if changes:
            fd = changes["file_delta"]
            dd = changes["dir_delta"]
            if fd or dd:
                print(f"  Δ files:{fd:+d} dirs:{dd:+d}", end="")
        print()

    return new_snap, changes


# ── Watch mode ────────────────────────────────────────────────
def watch(interval=3):
    print(f"SILK watching WORLDTREE every {interval}s  (Ctrl+C to stop)")
    print(f"  WORLDTREE: {WORLDTREE}")
    print(f"  MANIFEST:  {MANIFEST}")
    print()

    # Initial scan
    update(verbose=True)

    while True:
        time.sleep(interval)
        old_snap = load_ledger()

        # Quick hash: just total file count + total size
        # (faster than full re-scan for large trees)
        try:
            quick_files = sum(1 for _ in WORLDTREE.rglob("*") if _.is_file())
            quick_size  = sum(f.stat().st_size for f in WORLDTREE.rglob("*") if f.is_file())
        except Exception:
            continue

        old_files = old_snap.get("total_files", -1)
        old_size  = old_snap.get("total_bytes", -1)

        if quick_files != old_files or quick_size != old_size:
            print(f"  Change detected: {quick_files} files (was {old_files})")
            update(verbose=True)


# ── Status (quick, no write) ──────────────────────────────────
def status():
    snap = load_ledger()
    if not snap:
        print("No ledger found. Run: python3 silk.py")
        return
    ts    = snap.get("timestamp", "—")
    nf    = snap.get("total_files", 0)
    nd    = snap.get("total_dirs", 0)
    nb    = snap.get("total_bytes", 0)
    nc    = sum(b["chars"] for b in snap.get("blocks", {}).values())
    cov   = nc / UNICODE_17_TOTAL * 100
    print(f"WORLDTREE status (from ledger {ts})")
    print(f"  Files:    {nf:,}")
    print(f"  Folders:  {nd:,}")
    print(f"  Size:     {nb/1024/1024:.1f} MB")
    print(f"  Chars:    {nc:,} / {UNICODE_17_TOTAL:,}  ({cov:.1f}% of Unicode 17.0)")
    print(f"  MANIFEST: {MANIFEST}")


# ── CLI ───────────────────────────────────────────────────────
def main():
    p = argparse.ArgumentParser(description="SILK — WORLDTREE ledger")
    p.add_argument("--watch",    action="store_true", help="Watch mode")
    p.add_argument("--interval", type=int, default=3, help="Watch interval seconds")
    p.add_argument("--diff",     action="store_true", help="Show diff vs last scan")
    p.add_argument("--status",   action="store_true", help="Quick status")
    args = p.parse_args()

    SILK_DIR.mkdir(exist_ok=True)

    if args.status:
        status()
    elif args.watch:
        watch(args.interval)
    elif args.diff:
        old  = load_ledger()
        new  = scan_worldtree()
        diff = diff_snapshots(old, new)
        print(f"Files:  {diff['file_delta']:+d}")
        print(f"Dirs:   {diff['dir_delta']:+d}")
        print(f"Bytes:  {diff['bytes_delta']:+,}")
        if diff["added_files"]:
            print(f"\nAdded ({len(diff['added_files'])}):")
            for f in diff["added_files"][:30]: print(f"  + {f}")
        if diff["removed_files"]:
            print(f"\nRemoved ({len(diff['removed_files'])}):")
            for f in diff["removed_files"][:10]: print(f"  - {f}")
    else:
        update(verbose=True)
        print(f"  → {MANIFEST}")


if __name__ == "__main__":
    main()
