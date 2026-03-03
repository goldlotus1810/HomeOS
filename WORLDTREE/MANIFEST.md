# WORLDTREE — MANIFEST
> *Sổ cái tự động do SILK quản lý. Cập nhật mỗi khi có thay đổi trong WORLDTREE.*
> *Không chỉnh sửa tay file này.*

---

## SNAPSHOT

    Last updated:  2026-03-03 03:50 UTC
    Scan time:     2026-03-03T03:50:05.755018Z

    Total files:   141,547
    Total folders: 336
    Total size:    43.3 MB  (45,436,363 bytes)

## UNICODE COVERAGE

    Standard:      Unicode 17.0.0  (released 2025-09-09)
    Total 17.0:    159,801 characters
    In WORLDTREE:  141,079 characters
    Missing:       18,722 characters
    Coverage:      88.3%

    Missing mainly:
      - CJK Ext J (323B0–3347F): ~4,298  new in 17.0, needs UnicodeData.txt
      - CJK Ext B–H:             ~18,000 ideographs beyond Python 3.12 data
      - New 17.0 scripts:        Tai Yo, Beria Erfe (partial coverage)

    To reach 100%:
      wget https://unicode.org/Public/17.0.0/ucd/UnicodeData.txt
      python3 SILK/build_missing.py UnicodeData.txt

## LAST CHANGE

    Recorded:      2026-03-03 03:50 UTC
    Files Δ:       +0  (0 added, 0 removed)
    Folders Δ:     +0  (0 added, 0 removed)
    Size Δ:        +9,865 bytes

## WORLDTREE STRUCTURE

    01-Scripts/  (9 items)
      01-European/  (40 items)
      02-WestAsian/  (33 items)
      03-SouthAsian/  (51 items)
      04-SouthEastAsian/  (24 items)
      05-CentralAsian/  (10 items)
      06-EastAsian/  (6 items)
      07-African/  (18 items)
      08-American/  (7 items)
      09-Ancient/  (23 items)
    02-Symbols/  (19 items)
      01-Punctuation/  (6 items)
      02-Alphanumeric/  (7 items)
      03-Currency/  (34 items)
      04-Arrows/  (5 items)
      05-Mathematical/  (5 items)
      06-Technical/  (3 items)
      07-Geometric/  (4 items)
      08-Misc/  (3 items)
      09-Braille/  (257 items)
      10-Musical/  (4 items)
      11-Shorthand/  (2 items)
      12-SignWriting/  (673 items)
      13-ModifierTone/  (33 items)
      14-Roman/  (15 items)
      15-HalfFullwidth/  (226 items)
      ... and 4 more subdirectories
    03-Emoji/  (9 items)
      01-Faces/  (81 items)
      02-NatureObjects/  (769 items)
      03-Transport/  (119 items)
      04-Alchemical/  (125 items)
      05-Supplemental/  (257 items)
      06-Chess/  (99 items)
      07-ExtendedA/  (108 items)
      08-Games/  (3 items)
      09-EnclosedIdeographic/  (65 items)

## TOP BLOCKS BY CHARACTER COUNT

| Block                                             |   Chars | Folder |
|---------------------------------------------------|---------|--------|
| CJK Extension B                               |  42,720 | 01-Scripts/06-EastAsian/03-Chinese/07-CJK_ExtB |
| CJK Unified Ideographs                        |  20,992 | 01-Scripts/06-EastAsian/03-Chinese/05-CJK_Unified |
| Hangul Syllables                              |  11,172 | 01-Scripts/06-EastAsian/01-Korean/05-HangulSyllables |
| CJK Extension F                               |   7,473 | 01-Scripts/06-EastAsian/03-Chinese/11-CJK_ExtF |
| CJK Extension A                               |   6,592 | 01-Scripts/06-EastAsian/03-Chinese/06-CJK_ExtA |
| CJK Extension E                               |   5,762 | 01-Scripts/06-EastAsian/03-Chinese/10-CJK_ExtE |
| CJK Extension G                               |   4,939 | 01-Scripts/06-EastAsian/03-Chinese/12-CJK_ExtG |
| CJK Extension H                               |   4,192 | 01-Scripts/06-EastAsian/03-Chinese/13-CJK_ExtH |
| CJK Extension C                               |   4,154 | 01-Scripts/06-EastAsian/03-Chinese/08-CJK_ExtC |
| Yi Syllables                                  |   1,165 | 01-Scripts/06-EastAsian/04-Yi/01-YiSyllables |
| Egyptian Hieroglyphs                          |   1,072 | 01-Scripts/09-Ancient/18-EgyptianHieroglyphs |
| Mathematical Alphanumeric Symbols             |     996 | 02-Symbols/02-Alphanumeric/06-MathAlphanumeric |
| Cuneiform                                     |     922 | 01-Scripts/09-Ancient/14-Cuneiform |
| Tangut Components                             |     768 | 01-Scripts/06-EastAsian/03-Chinese/22-TangutComponents |
| Miscellaneous Symbols and Pictographs         |     768 | 03-Emoji/02-NatureObjects |
| Sutton SignWriting                            |     672 | 02-Symbols/12-SignWriting |
| Unified Canadian Aboriginal Syllabics         |     640 | 01-Scripts/08-American/05-UCAS |
| Arabic Presentation Forms-A                   |     631 | 01-Scripts/02-WestAsian/08-ArabicPresA |
| Anatolian Hieroglyphs                         |     583 | 01-Scripts/09-Ancient/21-AnatolianHieroglyphs |
| Bamum Supplement                              |     569 | 01-Scripts/07-African/09-BamumSup |
| CJK Compatibility Ideographs Supplement       |     542 | 01-Scripts/06-EastAsian/03-Chinese/17-CJKCompatSup |
| CJK Compatibility Ideographs                  |     472 | 01-Scripts/06-EastAsian/03-Chinese/16-CJKCompat |
| Khitan Small Script                           |     470 | 01-Scripts/06-EastAsian/03-Chinese/26-KhitanSmall |
| Nushu                                         |     396 | 01-Scripts/06-EastAsian/03-Chinese/25-Nushu |
| Ethiopic                                      |     358 | 01-Scripts/07-African/01-Ethiopic |
| Linear A                                      |     341 | 01-Scripts/09-Ancient/23-LinearA |
| Vai                                           |     300 | 01-Scripts/07-African/13-Vai |
| Latin Extended Additional                     |     256 | 01-Scripts/01-European/10-LatinExtAdditional |
| Cyrillic                                      |     256 | 01-Scripts/01-European/21-Cyrillic |
| Arabic                                        |     256 | 01-Scripts/02-WestAsian/03-Arabic |
| Hangul Jamo                                   |     256 | 01-Scripts/06-EastAsian/01-Korean/01-HangulJamo |
| Mathematical Operators                        |     256 | 02-Symbols/05-Mathematical/01-Operators |
| Supplemental Mathematical Operators           |     256 | 02-Symbols/05-Mathematical/05-SupplementalOps |
| Miscellaneous Technical                       |     256 | 02-Symbols/06-Technical/01-MiscTechnical |
| Miscellaneous Symbols                         |     256 | 02-Symbols/08-Misc/01-MiscSymbols |
| Braille Patterns                              |     256 | 02-Symbols/09-Braille |
| CJK Compatibility                             |     256 | 02-Symbols/19-CJKCompatibility |
| Supplemental Symbols and Pictographs          |     256 | 03-Emoji/05-Supplemental |
| Enclosed CJK Letters and Months               |     255 | 02-Symbols/18-EnclosedCJK |
| Miscellaneous Symbols and Arrows              |     253 | 02-Symbols/04-Arrows/04-MiscArrows |
| Byzantine Musical Symbols                     |     246 | 02-Symbols/10-Musical/01-Byzantine |
| Greek Extended                                |     233 | 01-Scripts/01-European/19-GreekExtended |
| Musical Symbols                               |     233 | 02-Symbols/10-Musical/02-Western |
| Halfwidth and Fullwidth Forms                 |     225 | 02-Symbols/15-HalfFullwidth |
| CJK Extension D                               |     222 | 01-Scripts/06-EastAsian/03-Chinese/09-CJK_ExtD |
| Kangxi Radicals                               |     214 | 01-Scripts/06-EastAsian/03-Chinese/03-KangxiRadicals |
| Mende Kikakui                                 |     213 | 01-Scripts/07-African/12-MendeKikakui |
| Symbols for Legacy Computing                  |     212 | 02-Symbols/16-LegacyComputing/01-Main |
| Tibetan                                       |     211 | 01-Scripts/05-CentralAsian/01-Tibetan |
| Latin Extended-B                              |     208 | 01-Scripts/01-European/04-LatinExtB |
| Enclosed Alphanumeric Supplement              |     200 | 02-Symbols/02-Alphanumeric/05-EnclosedAlphaSup |
| Early Dynastic Cuneiform                      |     196 | 01-Scripts/09-Ancient/16-EarlyDynasticCuneiform |
| Latin Extended-D                              |     193 | 01-Scripts/01-European/06-LatinExtD |
| Dingbats                                      |     192 | 02-Symbols/08-Misc/02-Dingbats |
| Znamenny Musical Notation                     |     185 | 02-Symbols/10-Musical/04-Znamenny |
| Myanmar                                       |     160 | 01-Scripts/04-SouthEastAsian/03-Myanmar |
| Enclosed Alphanumerics                        |     160 | 02-Symbols/02-Alphanumeric/04-EnclosedAlpha |
| Mongolian                                     |     158 | 01-Scripts/05-CentralAsian/02-Mongolian |
| Supplemental Arrows-C                         |     150 | 02-Symbols/04-Arrows/05-SupplementalC |
| Miao                                          |     149 | 01-Scripts/06-EastAsian/06-Miao |

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
