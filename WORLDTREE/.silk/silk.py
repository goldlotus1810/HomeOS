#!/usr/bin/env python3
"""
silk.py — Sổ cái tự động của WORLDTREE
========================================
Chạy lệnh này để:
  python3 .silk/silk.py          → đếm lại, cập nhật MANIFEST.md + ledger
  python3 .silk/silk.py diff     → so sánh với snapshot cũ, ghi changelog
  python3 .silk/silk.py watch    → theo dõi liên tục (polling mỗi 5 giây)

Nguyên tắc bất biến:
  - Không bao giờ xóa hay đổi tên file/folder trong WORLDTREE
  - Chỉ THÊM VÀO
  - Mọi thay đổi được ghi tự động vào ledger.json + MANIFEST.md
"""

import os, sys, json, hashlib, datetime, time

ROOT     = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SILK_DIR = os.path.join(ROOT, ".silk")
LEDGER   = os.path.join(SILK_DIR, "ledger.json")
SNAP     = os.path.join(SILK_DIR, "snapshot.json")
MANIFEST = os.path.join(ROOT, "MANIFEST.md")

UNICODE_17_TOTAL  = 159_801
UNICODE_17_BLOCKS = 338

# ─────────────────────────────────────────────────────────
def scan(root=ROOT):
    """Walk WORLDTREE, return stats dict."""
    dirs, char_files, other_files = [], [], []

    for r, ds, fs in os.walk(root):
        ds[:] = sorted(d for d in ds if d != ".silk")
        rel_dir = os.path.relpath(r, root)
        if rel_dir != ".":
            dirs.append(rel_dir)
        for f in sorted(fs):
            if f == "MANIFEST.md":
                continue
            rel = os.path.relpath(os.path.join(r, f), root)
            if f.startswith("U") and f.endswith(".md") and not f == "INDEX.md":
                char_files.append(rel)
            else:
                other_files.append(rel)

    return {
        "dirs":        dirs,
        "char_files":  char_files,
        "other_files": other_files,
    }


def hash_file(path):
    return hashlib.sha256(open(path, "rb").read()).hexdigest()[:16]


def make_snapshot(root=ROOT):
    snap = {}
    for r, ds, fs in os.walk(root):
        ds[:] = sorted(d for d in ds if d != ".silk")
        for f in sorted(fs):
            if f == "MANIFEST.md":
                continue
            p   = os.path.join(r, f)
            rel = os.path.relpath(p, root)
            snap[rel] = hash_file(p)
    return snap


def diff_snapshots(old_snap, new_snap):
    added   = [k for k in new_snap if k not in old_snap]
    removed = [k for k in old_snap if k not in new_snap]  # should never happen
    return added, removed


def block_stats(root=ROOT):
    """Per-block character count from folder structure."""
    stats = []
    for r, ds, fs in os.walk(root):
        ds[:] = sorted(d for d in ds if d != ".silk")
        rel = os.path.relpath(r, root)
        if rel == ".":
            continue
        idx_path = os.path.join(r, "INDEX.md")
        if not os.path.exists(idx_path):
            continue
        # Parse INDEX.md for range + assigned
        block_name, rng, assigned, slots = rel.split("/")[-1], "?", 0, 0
        char_count = sum(
            1 for f in fs
            if f.startswith("U") and f.endswith(".md") and f != "INDEX.md"
        )
        # Read range from INDEX.md
        try:
            with open(idx_path, encoding="utf-8") as fh:
                for line in fh:
                    if line.startswith("# "):
                        block_name = line[2:].strip()
                    elif line.startswith("range:"):
                        rng = line.split(":", 1)[1].strip()
                    elif line.startswith("slots:"):
                        slots = int(line.split(":", 1)[1].strip())
                    elif line.startswith("assigned:"):
                        assigned = int(line.split(":", 1)[1].strip())
        except:
            pass
        stats.append({
            "path":     rel,
            "name":     block_name,
            "range":    rng,
            "slots":    slots,
            "assigned": assigned,
            "files":    char_count,
        })
    return stats


# ─────────────────────────────────────────────────────────
def load_ledger():
    if os.path.exists(LEDGER):
        return json.load(open(LEDGER, encoding="utf-8"))
    return {
        "worldtree_version": "1.0.0",
        "unicode_version":   "17.0.0",
        "built_at":          now_iso(),
        "changelog":         [],
    }


def now_iso():
    return datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

def now_date():
    return datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d")


def save_ledger(data):
    os.makedirs(SILK_DIR, exist_ok=True)
    json.dump(data, open(LEDGER, "w", encoding="utf-8"), ensure_ascii=False, indent=2)


# ─────────────────────────────────────────────────────────
def build_manifest(stats_data, ledger_data):
    """Regenerate MANIFEST.md from current state."""
    total_chars  = sum(b["files"] for b in stats_data)
    total_dirs   = len(stats_data)
    missing      = UNICODE_17_TOTAL - total_chars
    pct          = total_chars / UNICODE_17_TOTAL * 100

    # Changelog table
    changelog_rows = ""
    for entry in reversed(ledger_data.get("changelog", [])):
        changelog_rows += (
            f"| {entry.get('version','?')} "
            f"| {entry.get('date','?')} "
            f"| {entry.get('unicode','17.0.0')} "
            f"| {entry.get('summary','?')} |\n"
        )

    # Top blocks
    top = sorted(stats_data, key=lambda x: x["files"], reverse=True)[:20]
    top_rows = ""
    for b in top:
        pct_b = b["files"] / b["slots"] * 100 if b["slots"] else 0
        top_rows += f"| {b['name']} | {b['range']} | {b['files']:,} | {b['slots']:,} | {pct_b:.0f}% |\n"

    # All blocks
    all_rows = ""
    for b in sorted(stats_data, key=lambda x: x["path"]):
        pct_b = b["files"] / b["slots"] * 100 if b["slots"] else 0
        flag  = "✓" if pct_b >= 99 else ("⚠" if pct_b > 0 else "✗")
        all_rows += f"| {b['name']} | {b['range']} | {b['files']:,} | {pct_b:.0f}% {flag} |\n"

    manifest = f"""\
# WORLDTREE / SILK — Sổ Cái Unicode
> Tự động tạo bởi `.silk/silk.py` · Cập nhật lần cuối: {now_iso()}

---

## PHIÊN BẢN

| | |
|---|---|
| Unicode Standard | **17.0.0** (2025-09-09) |
| Nguồn | unicode.org/versions/Unicode17.0.0/ |
| Build | {now_date()} |
| Ledger | `.silk/ledger.json` |
| Snapshot | `.silk/snapshot.json` |

---

## PHỦ SÓNG HIỆN TẠI

```
Unicode 17.0 tổng:    {UNICODE_17_TOTAL:>10,} ký tự
Có trong WORLDTREE:   {total_chars:>10,} ký tự
Thiếu:                {missing:>10,} ký tự
Phủ sóng:             {pct:>10.1f}%

Block Unicode 17.0:   {UNICODE_17_BLOCKS:>10}
Block trong cây:      {total_dirs:>10}
```

---

## THIẾU Ở ĐÂU?

Python 3.12 chỉ có Unicode data 15.0.0.
Các block mới 16.0/17.0 và một số CJK lớn chưa có data:

| Block | Range | ~Thiếu | Lý do |
|-------|-------|--------|-------|
| CJK Extension J | 323B0–3347F | 4,298 | 17.0 NEW |
| CJK Extension G | 30000–3134F | 4,939 | post-15.0 |
| CJK Extension H | 31350–323AF | 4,192 | post-15.0 |
| CJK Extension I | 2EBF0–2EE5F | 622 | post-15.0 |
| CJK Extension B (phần) | 20000–2A6DF | ~4k | partial |
| Garay | 10D40–10D8F | 64 | 16.0 new |
| Todhri | 105C0–105FF | 48 | 16.0 new |
| Kirat Rai | 16D40–16D7F | 64 | 16.0 new |
| Gurung Khema | 16100–1613F | 64 | 16.0 new |
| Ol Onal | 1E5D0–1E5FF | 48 | 16.0 new |
| Tai Yo | 1E6C0–1E6FF | 55 | 17.0 NEW |
| Sidetic | 10940–1095F | 26 | 17.0 NEW |
| Tolong Siki | 11DB0–11DEF | 54 | 17.0 NEW |
| Beria Erfe | 16EA0–16EDF | 50 | 17.0 NEW |
| Tangut Components Sup | 18D80–18DFF | 115 | 17.0 NEW |
| Misc Symbols Sup | 1CEC0–1CEFF | 34 | 17.0 NEW |
| Sharada Supplement | 11B60–11B7F | 8 | 17.0 NEW |

**Để thêm vào:**
```bash
# Thêm file thủ công vào đúng thư mục, rồi chạy:
python3 .silk/silk.py
# → tự động phát hiện, ghi vào ledger.json, tái tạo MANIFEST.md
```

---

## CẤU TRÚC CÂY

```
WORLDTREE/
├── MANIFEST.md                     ← file này (sổ cái hiển thị)
├── .silk/
│   ├── silk.py                     ← công cụ sổ cái (chạy để cập nhật)
│   ├── ledger.json                 ← lịch sử mọi thay đổi
│   └── snapshot.json               ← SHA-256 mọi file
│
├── 01-Scripts/                     Scripts — chữ viết loài người
│   ├── 01-European/                41 block (Latin, Greek, Cyrillic...)
│   ├── 02-WestAsian/               33 block (Hebrew, Arabic, Syriac...)
│   ├── 03-SouthAsian/              51 block (Devanagari, Tamil, Kawi...)
│   ├── 04-SouthEastAsian/          23 block (Thai, Khmer, Myanmar...)
│   ├── 05-CentralAsian/            9 block  (Tibetan, Mongolian...)
│   ├── 06-EastAsian/
│   │   ├── 01-Korean/              11,172 Hangul syllables + Jamo
│   │   ├── 02-Japanese/            Hiragana, Katakana
│   │   ├── 03-Chinese/             CJK A–J (87,000+ ideographs)
│   │   ├── 04-Yi/
│   │   └── 05-Lisu/
│   ├── 07-African/                 17 block (Ethiopic, Coptic, Adlam...)
│   ├── 08-American/                Cherokee, Deseret, UCAS
│   └── 09-Ancient/                 Cuneiform, Egyptian, Linear B...
│
├── 02-Symbols/                     Ký hiệu & Toán học
│   ├── 01-Punctuation/
│   ├── 02-Alphanumeric/
│   ├── 03-Currency/                33 ký hiệu tiền tệ
│   ├── 04-Arrows/                  5 block
│   ├── 05-Mathematical/            688+ toán tử
│   ├── 06-Technical/
│   ├── 07-Geometric/               Box Drawing, Geometric Shapes (U+25A0–25FF)
│   ├── 08-Misc/                    Dingbats, Misc Symbols
│   ├── 09-Braille/                 256 Braille patterns
│   ├── 10-Musical/                 4 block
│   ├── 11-Shorthand/
│   └── 14-LegacyComputing/
│
└── 03-Emoji/
    ├── 01-Faces/                   U+1F600–1F64F
    ├── 02-NatureObjects/
    ├── 03-Transport/
    ├── 04-Alchemical/
    ├── 05-Supplemental/
    ├── 06-Chess/
    ├── 07-ExtendedA/
    └── 08-Games/                   Mahjong, Domino, Cards
```

---

## FORMAT MỖI FILE LÁ

```
# ● — BLACK CIRCLE

codepoint: U+25CF
decimal:   9679
utf-8:     e2 97 8f
utf-32:    000025CF

category:  So (Other Symbol)
block:     Geometric Shapes  (U+25A0–U+25FF)
bidi:      ON

## Cross-references

→ U+2B24  ⬤  BLACK LARGE CIRCLE
→ U+1F311 🌑  NEW MOON SYMBOL
→ U+1F534 🔴  LARGE RED CIRCLE

## Character

●
```

---

## TOP 20 BLOCK LỚN NHẤT

| Block | Range | Chars | Slots | % |
|-------|-------|-------|-------|---|
{top_rows}
---

## TẤT CẢ BLOCK

| Block | Range | Chars | % |
|-------|-------|-------|---|
{all_rows}
---

## NGUYÊN TẮC BẤT BIẾN

```
1. Tên thư mục WORLDTREE    — không bao giờ đổi
2. Tên thư mục con          — không bao giờ đổi
3. Tên file                 — không bao giờ đổi
4. Chỉ THÊM VÀO             — không xóa, không đổi tên
5. Mọi thay đổi             → .silk/ledger.json + MANIFEST.md
```

---

## LỊCH SỬ THAY ĐỔI

| Version | Ngày | Unicode | Thay đổi |
|---------|------|---------|----------|
{changelog_rows}
*Chi tiết đầy đủ: `.silk/ledger.json`*

---
*WORLDTREE — nền tảng Unicode của SILK Knowledge Tree*
"""
    return manifest


# ─────────────────────────────────────────────────────────
def run_update(verbose=True):
    """Main: scan → diff → update ledger → write MANIFEST."""
    if verbose:
        print("🌳 WORLDTREE / SILK — scanning...")

    # 1. Load old snapshot
    old_snap = {}
    if os.path.exists(SNAP):
        old_snap = json.load(open(SNAP, encoding="utf-8"))

    # 2. New snapshot
    new_snap = make_snapshot(ROOT)

    # 3. Diff
    added, removed = diff_snapshots(old_snap, new_snap)

    # 4. Load ledger
    ledger = load_ledger()

    # 5. Record changes
    if added or not old_snap:
        n_chars = sum(1 for p in added if os.path.basename(p).startswith("U") and p.endswith(".md"))
        n_dirs  = sum(1 for p in added if "INDEX.md" in p)
        summary = (
            f"+{len(added)} file(s)"
            + (f", {n_chars} char(s)" if n_chars else "")
            + (f", {n_dirs} block(s)" if n_dirs else "")
            if added else
            f"Initial build: {len(new_snap)} files"
        )
        version = ledger.get("worldtree_version", "1.0.0")
        if added:
            # bump patch version
            parts = version.split(".")
            parts[-1] = str(int(parts[-1]) + 1)
            version = ".".join(parts)
            ledger["worldtree_version"] = version

        ledger["changelog"].append({
            "version":  version,
            "date":     now_date(),
            "at":       now_iso(),
            "unicode":  "17.0.0",
            "summary":  summary,
            "added":    added[:200],   # cap at 200 entries
            "removed":  removed,
        })
        if verbose and added:
            print(f"   + {len(added)} new file(s) detected")
        if removed and verbose:
            print(f"   ⚠ {len(removed)} file(s) missing (should not happen — immutability violation?)")

    # 6. Save snapshot
    os.makedirs(SILK_DIR, exist_ok=True)
    json.dump(new_snap, open(SNAP, "w", encoding="utf-8"), ensure_ascii=False, indent=2)

    # 7. Collect block stats
    stats = block_stats(ROOT)
    total_chars = sum(b["files"] for b in stats)
    total_dirs  = len(stats)

    # Update ledger totals
    ledger["built_at"]      = now_iso()
    ledger["total_chars"]   = total_chars
    ledger["total_blocks"]  = total_dirs
    ledger["total_files"]   = len(new_snap)
    save_ledger(ledger)

    # 8. Write MANIFEST.md
    manifest_text = build_manifest(stats, ledger)
    open(MANIFEST, "w", encoding="utf-8").write(manifest_text)

    if verbose:
        missing = UNICODE_17_TOTAL - total_chars
        pct     = total_chars / UNICODE_17_TOTAL * 100
        print(f"   ✓ {total_chars:,} chars  |  {total_dirs} blocks  |  {pct:.1f}% coverage")
        print(f"   ✗ {missing:,} chars missing (see MANIFEST.md)")
        print(f"   → MANIFEST.md updated")
        print(f"   → .silk/ledger.json updated  (v{ledger['worldtree_version']})")

    return total_chars, total_dirs


def run_watch(interval=5):
    """Poll for changes every N seconds."""
    print(f"👁  Watching WORLDTREE every {interval}s... (Ctrl+C to stop)")
    last_count = 0
    while True:
        try:
            count = sum(
                1 for r, d, fs in os.walk(ROOT)
                for f in fs
                if f.startswith("U") and f.endswith(".md")
            )
            if count != last_count:
                print(f"\n[{now_iso()}] Change detected ({last_count} → {count} chars)")
                run_update(verbose=True)
                last_count = count
            else:
                print(".", end="", flush=True)
            time.sleep(interval)
        except KeyboardInterrupt:
            print("\nStopped.")
            break


# ─────────────────────────────────────────────────────────
if __name__ == "__main__":
    cmd = sys.argv[1] if len(sys.argv) > 1 else "update"

    if cmd == "watch":
        interval = int(sys.argv[2]) if len(sys.argv) > 2 else 5
        run_watch(interval)
    elif cmd == "diff":
        old_snap = json.load(open(SNAP)) if os.path.exists(SNAP) else {}
        new_snap = make_snapshot(ROOT)
        added, removed = diff_snapshots(old_snap, new_snap)
        print(f"Added:   {len(added)}")
        print(f"Removed: {len(removed)}")
        for p in added[:50]:
            print(f"  + {p}")
        if len(added) > 50:
            print(f"  ... and {len(added)-50} more")
    else:
        run_update(verbose=True)
