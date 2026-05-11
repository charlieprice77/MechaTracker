# Mechabellum Overlay — Architecture Doc & Handoff

This is a handoff document for the Mechabellum FFA (Brawl) overlay/analytics project. It captures everything decided so far, the open questions, and the recommended path forward — enough that a fresh Claude Code instance (or yourself in a week) can pick it up cleanly.

---

## Project goal

A Mechabellum companion tool focused on 4-player FFA (Brawl). The first feature is an efficiency rating for **unit drop offers** shown each round, computed from known game constants so it works without any sample-size of games (unlike HS deck tracker's mulligan tool). Later features extend to cards (economy upgrades, unit improvements with synergy detection, counter detection for Improved Unit cards).

The tool should ideally be a **desktop overlay** in the spirit of HS deck tracker — looks nice, sits over the game while playing.

---

## Game economics (reference)

### Unit drops
- Each round (round 2+) the player is offered **4 unit drops** and **4 cards**.
- Drop format: `[qty]× R[level] [unit] for [cost], skip for [skip reward]`.
- The drop cost varies per offer; the skip reward is fixed by round.
- The player picks **one drop and one card per round** (or skips, which gives supply: +50 for cards, round-based for drops).

### Skip rewards by round
| Round | Skip reward |
|---|---|
| 2 | +100 |
| 3 | +150 |
| 4 | +200 |
| 5 | +250 |
| 6–10 | +300 |

(Cards skip = +50, all rounds.)

### Brawl sell value
- Selling a unit in Brawl returns **50%** of its full unit value (1v1/2v2 returns 100% — different mode).
- Unit value = `base recruitment cost + (level − 1) × upgrade cost`.
- Upgrade cost is flat across levels (lvl 8→9 costs the same as lvl 1→2). Levels go to 9.
- Unlock cost is paid **once** to enable recruiting a unit; it does **not** add to a unit's value (you can't sell back the unlock).

### Efficiency formulas (current calculator)

```
net cost     = drop cost + skip reward
unit value   = (base cost + (level − 1) × upgrade cost) × quantity
unlock saving = unlock cost, if not already unlocked, else 0

ratio        = (unit value + unlock saving) ÷ net cost
sale return  = (unit value × 0.5) − net cost      ← signed; can be negative
```

The ratio always factors in unlock saving (since taking the drop saves you the future unlock cost if you'd recruit more of that unit). Sale return is the pure "took it, sold it, what did I net" number — unlock saving excluded since you're not recruiting more.

### Colour scale
| Ratio | Colour | Label |
|---|---|---|
| < 1.50 | dark red | Poor |
| 1.50–1.74 | dark amber | Below average |
| 1.75–1.99 | amber | Decent |
| 2.00–2.24 | dark green | Good |
| 2.25–2.74 | green | Great |
| ≥ 2.75 | deep green | Exceptional |

---

## Unit database (current)

CSV the user provided is the source of truth. Columns: Unit Name, Recruitment Cost, Unlock Cost, Upgrade Cost, Note.

Key facts:
- Most light units cost 100/200/300, upgrade for 0.5× cost
- Most giants cost 400, unlock 200, upgrade 200
- Titans (Abyss, Death Knell, Mountain, War Factory): cost 800, unlock 400, upgrade 400
- Overlord: 500/250/250
- Typhoon: drop-only, no unlock cost
- Wasp, Phoenix, Phantom Ray, Rhino: 50 unlock
- Hacker: 100 unlock
- Farseer, Scorpion, Wraith: 50 unlock (corrected — was 150)

Wiki source: <https://wiki.mbxmas.com/> — has unit pages with cost, unlock cost, good against / bad against, technologies. Could be scraped programmatically to keep the database fresh across patches.

---

## XML replay format (decoded)

Files live at: `steam/steamapps/common/mechabellum/ProjectData/replays` (approximate), extension `.grbr`. The user must manually save via Esc → Save Replay — **the game does not auto-save replays.** This is the single biggest constraint on the architecture.

### Key finding: offerings ARE in the XML before selection

Round 4 replay saved before selection shows:
- `<reinforceItems>` populated with the round's cards and drops
- `<actionRecords />` empty for that round

This means: **a save taken the moment offerings appear contains the offered IDs, before any selection is made.**

### XML structure (top-level)

```
<BattleRecord>
  <playerRecords>
    <PlayerRecord>           ← one per player
      <id>...</id>
      <name>Charlie</name>
      <data>
        <unitDatas>           ← tech tree state for all units (IDs 1-31, 2001, 2002)
          <unitData>
            <id>N</id>
            <techs>
              <tech data="..." />   ← tech IDs, unknown encoding
            </techs>
            <unlockedTechs />
          </unitData>
          ...
        </unitDatas>
      </data>
      <playerRoundRecords>
        <PlayerRoundRecord>
          <round>N</round>
          <playerData>
            <reactorCore>...</reactorCore>
            <supply>...</supply>
            <preRoundFightResult>Win|Lose|Deuce</preRoundFightResult>
            <units>
              <NewUnitData>     ← every unit on the board
                <id>N</id>          ← unit type ID
                <Index>N</Index>    ← board slot
                <Level>N</Level>
                <SellSupply>N</SellSupply>
                <Position><x/><y/></Position>
                ...
              </NewUnitData>
              ...
            </units>
            <shop>
              <unlockedUnits><int>N</int>...</unlockedUnits>  ← per-round unlock state
              <lockedUnits><int>N</int>...</lockedUnits>
              <BuyCount>N</BuyCount>
              <UnlockCount>N</UnlockCount>
            </shop>
            <bluepints>...</bluepints>           ← active blueprint cards (typo in game)
            <officers>...</officers>             ← commander
            <commanderSkills>...</commanderSkills>
            <towerStrengthenLevels>...</towerStrengthenLevels>
          </playerData>
          <actionRecords>
            <MatchActionData xsi:type="PAD_ChooseReinforceItem">
              <ID>NNN</ID>           ← chosen card or drop ID; 0 = skipped
              <Index>N</Index>       ← 0-3 = slot chosen, -1 = skipped
            </MatchActionData>
            <MatchActionData xsi:type="PAD_BuyUnit">...</MatchActionData>
            <MatchActionData xsi:type="PAD_UnlockUnit">...</MatchActionData>
            <MatchActionData xsi:type="PAD_MoveUnit">...</MatchActionData>
            <MatchActionData xsi:type="PAD_UpgradeUnit">...</MatchActionData>
            <MatchActionData xsi:type="PAD_FinishDeploy">...</MatchActionData>
            <MatchActionData xsi:type="PAD_ChooseAdvanceTeam">...</MatchActionData>
            <MatchActionData xsi:type="PAD_ReleaseCommanderSkill">...</MatchActionData>
            ...
          </actionRecords>
        </PlayerRoundRecord>
        ...
      </playerRoundRecords>
    </PlayerRecord>
    ...
  </playerRecords>
  <matchDatas>
    <MatchSnapshotData>
      <round>N</round>
      <reinforceItems>           ← THE OFFERED ITEMS for this round
        <ArrayOfInt>             ← 4 card IDs
          <int>...</int>×4
        </ArrayOfInt>
        <ArrayOfInt>             ← 4 drop IDs
          <int>...</int>×4
        </ArrayOfInt>
      </reinforceItems>
      <lastFightResult>...</lastFightResult>
      ...
    </MatchSnapshotData>
    ...
  </matchDatas>
  <BattleInfo>
    <BattleID>YYYYMMDD--XXXXXXXXX</BattleID>
    <MatchMode>VS_4_Scuffle</MatchMode>  ← FFA / Brawl
    <GameMode>Normal</GameMode>
    <MaxRound>10</MaxRound>
    <ScoreMode>ReduceScore</ScoreMode>
    ...
  </BattleInfo>
</BattleRecord>
```

### Drop ID encoding (cracked)

Format: **`30 | round(1 digit) | qty(1 digit) | lvl(1 digit) | uid(variable, 1-2 digits)`**

Verified on 15 drops across 3 rounds. Examples:
- `302225` → round 2, 2× Lv2, uid=5 (Rhino)
- `3021123` → round 2, 1× Lv1, uid=23 (Sandworm)
- `3024130` → round 2, 4× Lv1, uid=30 (Void Eye)
- `3035131` → round 3, 5× Lv1, uid=31 (Vortex)
- `304211` → round 4, 2× Lv1, uid=1 (Fortress)
- `3042218` → round 4, 2× Lv2, uid=18 (Wraith)

### Unit internal ID map (14 of ~33 known)

| uid | Unit |
|---|---|
| 1 | Fortress |
| 3 | Vulcan |
| 5 | Rhino |
| 15 | Arclight |
| 16 | Phoenix |
| 18 | Wraith |
| 19 | Scorpion |
| 20 | Fire Badger |
| 22 | Typhoon |
| 23 | Sandworm |
| 25 | Phantom Ray |
| 26 | Farseer |
| 30 | Void Eye |
| 31 | Vortex |

Remaining ~19 IDs to discover through more replays. These IDs match the `<id>` field used in `<unitData>` and `<NewUnitData>`.

### Card ID encoding (NOT decoded yet)

Examples seen (with user-confirmed names):
| ID | Card |
|---|---|
| 30403 | Mass-Produced Melting Point |
| 31304 | Improved Sledgehammer |
| 30601 | Mass-Produced Wasp |
| 31101 | Fortified Overlord |
| 300001 | Missile Strike |
| 31001 | Subsidized Crawler |
| 20003 | Efficient Tech Research |
| 1200009 | Scorpion Assault |
| 1100001 | (unknown — Intensive Training, chosen R2 game 1) |
| 13030010 | Dominion Core |
| 31802 | Improved Wraith |
| 13030006 | Super Heavy Armor |
| 20002 | Advanced Offensive Tactics |
| 3022220 | (this is a drop ID, not a card — listed by mistake earlier) |

No obvious pattern yet. Cards probably have arbitrary IDs from a game database. Cracking these likely requires either many more examples or finding game data files (e.g. in Steam install).

Skip card: `ID=0`, `Index=-1` in the `PAD_ChooseReinforceItem` action.

---

## Architecture options (the decision)

The "manual save" constraint is the central architectural pivot. Three viable paths:

### Option A: Macro-triggered XML parsing

**Flow**: AutoHotkey (or similar) macro hits Esc → Save Replay → Esc when triggered. A file watcher reads the latest `.grbr`, parses the most recent round's `reinforceItems`, and shows the efficiency overlay.

**Triggering**: User hotkey, or a poll on the game window state (e.g. detecting the drop screen via window title or pixel sampling).

**Pros**:
- Data is reliable and complete (full board state, unlock state, tech state — not just the current offer)
- No OCR fragility
- Enables features beyond drops (board state analysis, synergy detection, counter detection) since the full game state is available
- Replay analysis features come for free with the same parser

**Cons**:
- Intrusive: forces a menu open mid-round every round (Esc → menu → Save Replay → Esc) — even with a macro, this animation/navigation is visible to the player
- Risk: if the game state changes (e.g. timer behaviour, menu structure) the macro breaks
- Risk: forcing the menu mid-round may cause issues in actual online play (multiplayer doesn't pause)

**Open question**: does Save Replay work mid-round in online multiplayer without disrupting the game? In the test files this was offline vs bots.

### Option B: Screen reading (OCR / template matching)

**Flow**: Continuously screenshot a region of the game window. Detect the drop UI. Run OCR or template-match unit icons to extract `[qty]× R[level] [unit] for [cost]`. Show overlay.

**Pros**:
- No game interaction, no menu navigation — completely passive
- Works in online multiplayer without disruption
- Easier to make latency-free

**Cons**:
- Fragile to UI changes, resolution, game updates, locale (non-English clients), accessibility settings
- OCR accuracy on unit names / numbers is workable but not 100%
- Doesn't get board state, tech state, unlock state — limited to whatever's on screen
- Limits feature scope to "things visible on the drop offer screen" — can't do counter detection (need to see other players' boards), can't reliably do synergy detection (need to know what's on user's board + active tech)

### Option C: Manual input (the current calculator)

**Flow**: Player taps in the offering. Overlay shows efficiency.

**Pros**:
- Zero technical risk; already working
- No dependency on game files or screen state
- Works in any client / any patch

**Cons**:
- Manual input every round is friction — defeats the "glanceable overlay" ideal
- Doesn't scale to the card features (synergy/counter detection) because those require knowing the broader game state

### Option D: Hybrid (recommended)

Three layers, each independently useful:

1. **Live efficiency calculator (manual input)** — what's already built. Works in any context, is a useful standalone tool, and is the foundation. Could be made into a small always-on-top window.

2. **Replay analyser (XML parser)** — post-game tool. Reads saved replays, decodes drops/cards/board state, shows what choices were made, what their efficiency was, what was offered, etc. Useful for review and content creation. **Builds the unit/card ID maps over time naturally.**

3. **Live overlay (XML + optional macro)** — if the user *opts in* to using a Save Replay macro, the same XML parser powers a live overlay. Falls back to manual input if not.

This way the project has value at every milestone, and the decision on the live-overlay-via-XML path can be deferred or revisited.

---

## Recommended next steps (for Claude Code or fresh instance)

In order of priority:

### 1. Python XML parser library
- Pure Python, no external deps beyond `xml.etree.ElementTree`.
- Load a `.grbr` file (note: it has a binary header before the XML — strip it or use a tolerant parser).
- Decode `reinforceItems` per round into structured data: `{round, cards: [...], drops: [{qty, level, unit_uid, raw_id, cost, ...}]}`.
- Decode `actionRecords` per round into chosen actions.
- Expose simple API: `parse_replay(path) → Replay` with rounds, players, etc.

This is the foundation. Get this solid before anything else.

### 2. Unit ID mapping completion
- The 14 known units (above) will grow as the parser is used.
- Build a mechanism to dump "unknown uid: N" warnings during parsing so unknown IDs surface immediately.
- Consider scraping the wiki to bootstrap if a more programmatic mapping is found.

### 3. CLI tool
- `mechabellum-analyse replay.grbr` prints a per-round breakdown of offered cards/drops, what was chosen, and the efficiency rating of each drop (using the formulas above).
- This is the replay analyser MVP. Useful by itself.

### 4. Card ID decoding (research)
- Look in the Steam install dir for game data files (.dat, .json, .csv, .bytes — Unity asset bundles maybe).
- If a card ID → name table exists in-game, it's the cleanest source.
- Otherwise, build the mapping from observed (id, user-supplied name) pairs over time.

### 5. Calculator GUI
- Port the web-based efficiency calculator to a small always-on-top desktop window.
- Python with PyQt6 / Tkinter / PyWebView (PyWebView lets you reuse the existing HTML/JS calculator).
- Add a "load latest replay" button that pre-fills the current round from a saved XML.

### 6. Live overlay (optional, deferred)
- Only build if the user decides the macro path is acceptable in their actual playstyle.
- File watcher on the replay folder, parses on change, updates overlay.
- AutoHotkey script to trigger Save Replay (separate from the Python tool — user installs and configures themselves).

### Tech stack recommendation
- **Python 3.11+** for parser, CLI, analyser. User is comfortable with Python + has C/firmware background.
- **PyWebView or PyQt6** for the desktop window — reuses the calculator's HTML/CSS or rewrites in native widgets. PyWebView is easier (just point it at a local HTML file); PyQt6 is more native-feeling.
- Avoid Electron unless the user really wants JS — Python is closer to user's comfort zone and parsing XML/files is more natural there.

---

## Open questions / known unknowns

1. **Does Save Replay work mid-round in online multiplayer?** Test files were vs bots. If it does, the macro approach is viable.
2. **Are drop offerings the same for all players, or per-player?** Affects whether a single replay decode tells us only the user's options or also opponents'. The `reinforceItems` is inside `<MatchSnapshotData>` (match-level), not inside `<PlayerRecord>`, which suggests global — but **verify this**.
3. **Tech ID encoding.** The `<tech data="N">` numbers (like 10401, 1201, 3001) likely encode which tech upgrades are available/active per unit. Decoding these unlocks the synergy detection feature later.
4. **Match mode detection.** `<MatchMode>VS_4_Scuffle</MatchMode>` is Brawl; need to know the strings for 1v1 / 2v2 / Survival so the tool can apply the correct sell ratio (50% Brawl vs 100% other modes).
5. **Card ID structure.** Open research problem.
6. **Game patches changing IDs.** Unit/card IDs may shift between patches. The tool needs a way to handle unknowns gracefully and a way to update the mapping.

---

## What's already built

A web-based efficiency calculator (HTML/JS) is iterated and working. Latest version computes:
- Single efficiency ratio with colour scale
- Sale return (signed)
- Full breakdown of inputs
- Round-based skip reward LUT
- Handles already-unlocked units, drop-only units (Typhoon), and quantity scaling

This can be packaged as the foundation of the GUI in step 5 above.

---

## Files referenced

- `Units_-_Sheet1__1_.csv` (user-provided) — unit database with corrected Farseer/Scorpion/Wraith unlocks
- Two `.grbr` replay files from custom FFA with bots, used to crack the drop ID format
