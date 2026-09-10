
# CITY OF LIES
## Master Game Bible + Scenario Script + Go Implementation Specification

**Document status:** V1 implementation blueprint  
**Backend:** Go  
**Frontend:** Web / Next.js / React / React Three Fiber  
**Primary mode:** Single-player real-time multi-agent investigation  
**Scenario:** Riverside Factory Incident  
**Primary design rule:** The LLM can express a character; the simulation engine owns reality.

---

> **How to use this document**
>
> This file is intentionally much more detailed than a normal README. Treat it as a hybrid of:
>
> 1. game design document,
> 2. screenplay/scenario bible,
> 3. multi-agent simulation specification,
> 4. backend architecture specification,
> 5. frontend implementation specification,
> 6. data contract,
> 7. test plan,
> 8. execution prompt for a coding agent.
>
> A coding agent must not “fill in the blanks” by inventing new rules where this file already defines them. If implementation constraints require a change, record the change as an ADR before modifying the design.

---

# 0. LOCKED DECISIONS

These decisions are considered **locked for V1** unless the repository owner explicitly changes them.

- Project name: **City of Lies**.
- V1 is a **browser game**, not a native desktop game.
- V1 uses a **3D isometric / orthographic presentation**.
- Backend is **Go**, not .NET, Node.js or Python.
- V1 has **one playable scenario** and **20 autonomous NPC agents**.
- The player is an independent investigator / journalist / fact-checker.
- The player wins by discovering the correct atomic Ground Truth and submitting it with enough evidence.
- The player loses when the primary false narrative becomes socially dominant.
- Default defeat threshold: **75% of eligible NPCs adopt the primary false narrative**.
- Default win evidence threshold: **70% Evidence Strength**.
- Default belief-adoption threshold for an individual NPC: **65% confidence**.
- NPCs are never omniscient.
- Ground Truth is never sent to an LLM.
- An LLM never writes directly to the database.
- An LLM never calculates win/loss.
- An LLM never mutates beliefs directly.
- The simulation must remain playable with LLM integration disabled.
- Multiplayer, procedural scenarios, voice, combat, crafting and open-world expansion are out of scope for V1.
- The backend remains a modular monolith for V1.
- Do not add Kafka, RabbitMQ, Redis or Kubernetes simply to make the architecture look “senior”.
- Reproducible simulation tests are more important than clever AI prompts.


# 1. PRODUCT THESIS

## 1.1 One-sentence pitch

**City of Lies is a real-time investigation game where the player must reconstruct the truth while misinformation spreads autonomously through a society of AI-driven characters with incomplete knowledge, individual beliefs, relationships and motives.**

## 1.2 The question the game explores

> If twenty autonomous people each see only part of an event, trust different people, remember things imperfectly and communicate under social pressure, how quickly can a believable lie become the accepted truth?

The game is not trying to lecture the player about “fake news”. The player should **feel** the mechanics of information distortion by making decisions under time pressure.

## 1.3 Player fantasy

The player fantasy is:

> “I am the only person trying to reconstruct the complete picture. Everyone around me knows something, believes something, wants something, and the city keeps talking even when I am not there.”

The player should never feel like a god who can open an admin panel and inspect every NPC's real belief.

## 1.4 Core tension

There are two races happening at the same time:

```text
INVESTIGATION RACE
find source → interview → inspect → corroborate → build evidence → submit truth

SOCIAL PROPAGATION RACE
observation → retelling → distortion → amplification → adoption → social dominance
```

Every minute spent investigating one lead gives the rumor network more time to evolve.

---

# 2. GAME OBJECTIVE AND OUTCOME RULES

## 2.1 Win

A session is won only when all of the following are true:

1. The player submits the correct answers for all required Ground Truth fields.
2. Evidence Strength is at least `truth_evidence_threshold`.
3. The primary false narrative ratio is below `defeat_false_belief_ratio`.
4. The session has not already reached a terminal loss state.

Default:

```text
truth_evidence_threshold      = 0.70
defeat_false_belief_ratio    = 0.75
belief_adoption_threshold    = 0.65
```

## 2.2 Loss

The session is immediately lost when:

```text
primary_false_narrative_ratio >= 0.75
```

Loss is evaluated after every belief-changing operation and after every public correction.

## 2.3 Why 75%

75% is used because it creates a readable “point of social dominance”. It is not presented as a scientific claim about real societies. It is simply a gameplay parameter.

## 2.4 Important distinction: claim vs narrative

A single NPC may believe one false claim but not the complete false narrative.

The primary false narrative is composed of three core false claims:

```text
F1: a chemical explosion happened.
F2: multiple workers died.
F3: the company is hiding those deaths.
```

For V1, an NPC adopts the **primary false narrative** when the weighted narrative score reaches the individual adoption threshold.

Suggested weights:

```text
chemical explosion   0.35
multiple deaths      0.40
company cover-up     0.25
```

Narrative score:

```text
score =
  belief(F1) * 0.35 +
  belief(F2) * 0.40 +
  belief(F3) * 0.25
```

Then:

```text
adopted = score >= 0.65
```

These weights are scenario-configurable.

---

# 3. V1 PLAYER EXPERIENCE

## 3.1 Target session length

Target first-play session:

```text
12–20 real minutes
```

The engine uses game time, but initial tuning should keep game time close enough to real time that the player experiences urgency.

## 3.2 Player verbs

Only these verbs are first-class in V1:

- **Interview** a person.
- **Inspect** a location.
- **Inspect / acquire** evidence.
- **Review** notebook.
- **Publish Correction**.
- **Submit Ground Truth**.

Do not add arbitrary RPG actions before these six feel good.

## 3.3 Player information surfaces

The player always sees:

- game time,
- global false-belief percentage,
- evidence strength,
- player credibility,
- recent public events.

The player does **not** automatically see:

- exact private belief score for each NPC,
- private memories,
- hidden motives,
- Ground Truth classifications,
- internal propagation coefficients,
- who will speak next,
- developer-only scenario tags.

## 3.4 Expected emotional arc

A successful session should roughly produce this sequence:

1. **Curiosity** — “What actually happened?”
2. **Confidence** — “This person probably knows.”
3. **Doubt** — “Their story conflicts with someone else.”
4. **Pressure** — “The false-belief meter is climbing.”
5. **Prioritization** — “I cannot interview everyone.”
6. **Risk** — “Do I publish now with incomplete evidence?”
7. **Resolution** — “I can prove the truth,” or “I waited too long.”

---

# 4. WORLD AND PRESENTATION

## 4.1 District

V1 takes place in **Riverside District**, a fictional compact urban district with subtle Vietnamese visual influences.

It must not replicate a real city street one-to-one.

## 4.2 Locations

Use eight gameplay locations:

1. `factory_gate` — Riverside Factory Gate
2. `factory_yard` — Factory Yard / evacuation area
3. `factory_office` — Management + maintenance office
4. `factory_storage` — affected storage building
5. `fire_station` — local fire station
6. `hospital` — clinic/hospital intake
7. `market` — local market / social hub
8. `cafe` — cafe / rumor hub
9. `news_office` — small local newsroom/blog office
10. `public_square` — public square
11. `residential` — residential block

The visual map may group some of these into 6–8 rendered building clusters, but the simulation location IDs remain distinct.

## 4.3 Camera

- Orthographic.
- Isometric default.
- Pan.
- Zoom.
- Limited rotation.
- No first-person.
- No free-flight debug controls in production UI.

Suggested default camera:

```text
position: [12, 16, 12]
target:   [0, 0, 0]
zoom:     tuned so 70–80% of district is readable
```

## 4.4 Visual priority

The 3D scene exists to answer:

- where are the characters?
- where is evidence?
- which places are becoming socially active?
- which characters just communicated?
- where should I investigate next?

Do not trade readability for visual effects.

---

# 5. SCENARIO 01 — RIVERSIDE FACTORY INCIDENT

## 5.1 Ground Truth summary

At approximately 17:58, an electrical fault in distribution panel DB-4 causes arcing and a localized fire in the storage building of Riverside Factory.

There is no major chemical explosion.

There are no fatalities.

Four people receive minor medical treatment.

A maintenance ticket had already reported overheating in DB-4 and had remained unresolved for twelve days.

Management is genuinely vulnerable to criticism for delayed maintenance, which makes the later cover-up rumor psychologically believable even though the alleged deaths do not exist.

## 5.2 Ground Truth atomic fields

The player's final submission must answer:

```text
event_type        = fire
origin            = electrical_distribution_panel_db4
cause_category    = electrical_fault
major_explosion   = false
chemical_release  = false
fatalities        = 0
injury_severity   = minor
management_issue  = delayed_maintenance
```

Required scoring fields for V1:

```text
event_type
cause_category
major_explosion
fatalities
```

Optional bonus truth fields:

```text
chemical_release
management_issue
origin
```

A player can win without identifying every bonus field, but the after-action report scores completeness.

## 5.3 Why the lie is believable

The lie must not be absurd.

Several truthful details make it fertile:

- witnesses heard a sharp “bang” from the electrical panel,
- heavy dark smoke appeared,
- ambulances arrived,
- management initially gave an evasive statement,
- a maintenance problem really had been ignored,
- a dramatic photo crop makes smoke appear much larger,
- the public cannot immediately access hospital or fire-service records.

This creates the core design principle:

> A good false narrative is built from fragments that are partly true.

## 5.4 Primary false narrative

```text
“A chemical explosion happened at Riverside Factory.
Several workers died, and management is hiding the real casualty count.”
```

## 5.5 Secondary rumors

Secondary rumors can exist but do not directly determine defeat:

- “A toxic cloud is moving toward the residential block.”
- “The factory was storing banned chemicals.”
- “The fire brigade was ordered not to speak.”
- “The hospital has been told to hide victims.”
- “The blogger has a confidential source confirming deaths.”

Secondary rumors create noise and force prioritization.

---

# 6. MASTER EVENT TIMELINE

This timeline is the **authoritative scenario seed**, not a hard-scripted sequence that must always play identically. The simulation scheduler may vary NPC-to-NPC propagation using a deterministic RNG seed, but the world event anchors below are fixed.

## 6.1 Pre-incident

### 17:40
Technician Đức checks a previous ticket for DB-4. Status remains unresolved.

### 17:47
Delivery driver Phúc passes the factory and sees a technical vehicle near the storage building.

### 17:52
Lan begins moving boxes near the storage area.

### 17:55
Minh performs a gate/security round.

## 6.2 Incident

### 17:58:10
DB-4 experiences electrical arcing.

### 17:58:13
Lan sees sparks.

### 17:58:15
A sharp electrical bang occurs.

This is deliberately audible enough for witnesses to later call it an “explosion”, but it is not a major blast.

### 17:58:25
Smoke begins building.

### 17:59
Local evacuation starts.

### 18:00
Security calls emergency services.

### 18:01
Tùng begins taking photos from outside the gate.

### 18:02
Hương sees workers coming out coughing.

### 18:03
First fire unit arrives.

### 18:04
Ambulance arrives.

### 18:05
Rumor seed 1 occurs:
Lan tells Hùng:

> “Tủ điện tóe lửa rồi nổ một cái, khói kín luôn.”

Hùng does not preserve the qualifier “tủ điện”.

His memory normalizes this to:

> “Lan nói bên trong có nổ.”

### 18:06
Hùng tells Thảo:

> “Hình như trong nhà máy có vụ nổ, xe cấp cứu vào rồi.”

This becomes the first explicit `claim_explosion_occurred` transmission.

### 18:07
Thảo repeats this to Hương and Phúc.

### 18:08
Tùng sends a dramatic smoke photo to Khoa.

### 18:09
Khoa receives Hùng's message:
> “Nghe nói có người bị kẹt trong đó.”

This statement is ambiguous and must not automatically create `multiple deaths`.

### 18:10
Khoa publishes an initial post:

> “Riverside Factory xảy ra sự cố lớn, nhân chứng nói nghe tiếng nổ. Xe cứu hỏa và cấp cứu đang có mặt.”

This post is mostly defensible.

### 18:12
A commenter/NPC paraphrase introduces:
> “Có thể là nổ hóa chất.”

Khoa updates headline with “nghi nổ hóa chất” without verified source.

### 18:13
Vy sees Khoa's post and shares it.

### 18:14
False-belief ratio is expected to cross approximately 15–25% depending on seed.

### 18:15
Hospital has already confirmed internally: four minor cases, zero fatalities.

Player can discover this but the city does not automatically know it.

### 18:16
Sơn, the manager, gives a poor public statement:

> “Chúng tôi đang kiểm tra và chưa thể bình luận về nguyên nhân hoặc thương vong.”

Although procedurally understandable, this increases cover-up suspicion.

### 18:18
A distorted chain reaches Vy:

> “Có người nói bệnh viện nhận nhiều công nhân và công ty đang giấu số liệu.”

Vy says in a live-style post:

> “Mình chưa xác nhận được nhưng có thông tin nói có người chết.”

This materially increases `claim_multiple_deaths`.

### 18:20
Expected false-belief range: 30–45%.

### 18:22
Khoa edits headline again:

> “Nghi vấn nổ hóa chất tại Riverside: nguồn tin nói có thương vong.”

He still has no verified death source.

### 18:24
Expected false-belief range: 40–55%.

### 18:25
Journalist Ngọc begins independent verification.

If the player reaches her with useful evidence, she can become a strong correction amplifier.

### 18:28
Community leader Liên requests official information.

### 18:30
If the player has done nothing, rumor pressure should become visibly critical.

Expected false-belief range: 55–65%.

### 18:32
Secondary toxic-cloud rumor may spawn if:
- chemical-explosion belief is high enough, and
- no credible public correction has been published.

### 18:34
Thầy Bình pushes back against toxic-cloud claims if he has access to air-sensor evidence or a credible player correction.

### 18:35+
From here, the scenario becomes highly state-dependent.

The scheduler must continue communication until:
- player wins,
- belief reaches defeat threshold,
- or optional hard session time limit is reached.

## 6.3 Optional hard limit

V1 may use:

```text
max_game_minutes = 50
```

If reached without a correct submission, treat as loss by “investigation deadline”.

Keep this disabled initially until playtesting proves it is needed.

---

# 7. NPC DESIGN PRINCIPLES

1. No NPC is a truth oracle.
2. Reliable NPCs can still be incomplete.
3. Unreliable NPCs can still possess important observations.
4. An NPC can intentionally evade without inventing facts.
5. Only NPCs with a sufficient deception trait plus a motive can intentionally lie.
6. “Lying” means contradicting something the NPC internally knows/believes for a motive; it does not grant hidden Ground Truth.
7. Memory provenance is retained.
8. Direct observation is stronger than hearsay.
9. Repetition can increase confidence, but repeated claims from the same source family should have diminishing impact.
10. Social influence and epistemic credibility are separate attributes.


# 8. NPC CAST — ALL 20 AGENTS
The following values are scenario seed values. They must be loaded from data files rather than embedded in Go source.
| ID | Name | Role | Influence | Credibility | Skepticism | Sociality | Deception | Receptiveness | Start |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| `agent_security_guard` | Minh Trần | Bảo vệ nhà máy | 0.38 | 0.78 | 0.72 | 0.35 | 0.10 | 0.38 | `factory_gate` |
| `agent_factory_worker_a` | Lan Phạm | Công nhân kho | 0.32 | 0.66 | 0.45 | 0.70 | 0.08 | 0.63 | `factory_yard` |
| `agent_factory_worker_b` | Hùng Nguyễn | Công nhân đóng gói | 0.44 | 0.52 | 0.28 | 0.88 | 0.18 | 0.82 | `cafe` |
| `agent_firefighter` | Quang Lê | Lính cứu hỏa | 0.55 | 0.95 | 0.90 | 0.30 | 0.02 | 0.20 | `fire_station` |
| `agent_paramedic` | Mai Vũ | Nhân viên cấp cứu | 0.42 | 0.90 | 0.86 | 0.35 | 0.01 | 0.25 | `hospital` |
| `agent_shop_owner` | Thảo Đỗ | Chủ quán tạp hóa | 0.58 | 0.62 | 0.40 | 0.90 | 0.12 | 0.68 | `market` |
| `agent_resident` | Bác Dũng | Cư dân khu Riverside | 0.36 | 0.58 | 0.50 | 0.62 | 0.04 | 0.57 | `residential` |
| `agent_blogger` | Khoa Bùi | Blogger địa phương | 0.82 | 0.48 | 0.33 | 0.80 | 0.35 | 0.65 | `news_office` |
| `agent_influencer` | Vy Hoàng | Influencer khu vực | 0.96 | 0.57 | 0.25 | 0.95 | 0.22 | 0.73 | `public_square` |
| `agent_taxi_driver` | Tuấn Đặng | Tài xế taxi | 0.43 | 0.51 | 0.38 | 0.84 | 0.08 | 0.71 | `public_square` |
| `agent_student` | An Nguyễn | Sinh viên | 0.47 | 0.46 | 0.57 | 0.86 | 0.03 | 0.76 | `cafe` |
| `agent_journalist` | Ngọc Trịnh | Phóng viên địa phương | 0.78 | 0.88 | 0.92 | 0.55 | 0.02 | 0.24 | `news_office` |
| `agent_factory_manager` | Sơn Vương | Quản lý nhà máy | 0.68 | 0.55 | 0.70 | 0.42 | 0.45 | 0.28 | `factory_office` |
| `agent_technician` | Đức Phan | Kỹ thuật viên bảo trì | 0.34 | 0.86 | 0.84 | 0.28 | 0.16 | 0.30 | `factory_office` |
| `agent_nurse` | Hà Lương | Y tá | 0.39 | 0.84 | 0.78 | 0.48 | 0.01 | 0.31 | `hospital` |
| `agent_vendor` | Cô Hương | Người bán hàng rong | 0.49 | 0.50 | 0.32 | 0.93 | 0.06 | 0.78 | `factory_gate` |
| `agent_retired_teacher` | Thầy Bình | Giáo viên nghỉ hưu | 0.61 | 0.79 | 0.82 | 0.50 | 0.00 | 0.29 | `residential` |
| `agent_delivery_driver` | Phúc Trần | Tài xế giao hàng | 0.41 | 0.49 | 0.36 | 0.72 | 0.05 | 0.73 | `market` |
| `agent_community_leader` | Cô Liên | Trưởng nhóm cộng đồng | 0.88 | 0.82 | 0.69 | 0.82 | 0.03 | 0.44 | `public_square` |
| `agent_photographer` | Tùng Lâm | Nhiếp ảnh gia tự do | 0.66 | 0.63 | 0.54 | 0.67 | 0.10 | 0.58 | `factory_gate` |

## 8.1 Individual character bible

### 8.2 Minh Trần — Bảo vệ nhà máy

**Agent ID:** `agent_security_guard`  
**Starting location:** `factory_gate`  
**Personality:** thận trọng, ít nói, quan sát tốt  
**Private motive:** Muốn bảo vệ công việc nhưng không muốn bao che cho tai nạn.

**Directly knows / witnessed**
- Thấy khói bốc từ kho điện; nghe tiếng nổ nhỏ kiểu aptomat/thiết bị chập, không phải nổ hóa chất; thấy 2 công nhân ho sặc khi chạy ra.

**Does not know**
- Không vào bên trong kho; không biết báo cáo bệnh viện; không biết nguyên nhân kỹ thuật cuối cùng.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.38,
  "credibility": 0.78,
  "skepticism": 0.72,
  "sociality": 0.35,
  "deception_tendency": 0.10,
  "receptiveness": 0.38
}
```

### 8.3 Lan Phạm — Công nhân kho

**Agent ID:** `agent_factory_worker_a`  
**Starting location:** `factory_yard`  
**Personality:** hoảng hốt, nói nhanh, dễ bị ảnh hưởng  
**Private motive:** Muốn mọi người biết điều kiện làm việc kém nhưng không cố tình bịa chuyện.

**Directly knows / witnessed**
- Ở gần kho khi cháy; thấy tia lửa ở tủ điện; bị ho do khói; được sơ cứu, không thấy ai tử vong.

**Does not know**
- Không biết số người bị thương toàn bộ; không biết lịch bảo trì.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.32,
  "credibility": 0.66,
  "skepticism": 0.45,
  "sociality": 0.70,
  "deception_tendency": 0.08,
  "receptiveness": 0.63
}
```

### 8.4 Hùng Nguyễn — Công nhân đóng gói

**Agent ID:** `agent_factory_worker_b`  
**Starting location:** `cafe`  
**Personality:** nhiều chuyện, thích kể chuyện, phóng đại vô thức  
**Private motive:** Muốn mình là người có tin nóng.

**Directly knows / witnessed**
- Không chứng kiến trực tiếp; nghe Lan nói 'có tiếng nổ ở tủ điện' và biến thành 'có vụ nổ'.

**Does not know**
- Hầu như toàn bộ sự thật kỹ thuật.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.44,
  "credibility": 0.52,
  "skepticism": 0.28,
  "sociality": 0.88,
  "deception_tendency": 0.18,
  "receptiveness": 0.82
}
```

### 8.5 Quang Lê — Lính cứu hỏa

**Agent ID:** `agent_firefighter`  
**Starting location:** `fire_station`  
**Personality:** điềm tĩnh, chuyên nghiệp, chỉ nói điều chắc chắn  
**Private motive:** Ưu tiên an toàn và thông tin chính xác.

**Directly knows / witnessed**
- Đội của anh dập cháy tại kho; không ghi nhận dấu vết nổ hóa chất; kết luận sơ bộ điểm cháy gần tủ phân phối điện.

**Does not know**
- Không có quyền kết luận pháp y cuối cùng; không biết động cơ của người lan tin.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.55,
  "credibility": 0.95,
  "skepticism": 0.90,
  "sociality": 0.30,
  "deception_tendency": 0.02,
  "receptiveness": 0.20
}
```

### 8.6 Mai Vũ — Nhân viên cấp cứu

**Agent ID:** `agent_paramedic`  
**Starting location:** `hospital`  
**Personality:** thực tế, rõ ràng, ưu tiên dữ liệu  
**Private motive:** Không muốn bệnh nhân bị khai thác làm tin giật gân.

**Directly knows / witnessed**
- Tiếp nhận 4 ca nhẹ: 2 hít khói, 1 trầy xước khi sơ tán, 1 hoảng loạn; không có tử vong.

**Does not know**
- Không biết chính xác nguồn cháy.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.42,
  "credibility": 0.90,
  "skepticism": 0.86,
  "sociality": 0.35,
  "deception_tendency": 0.01,
  "receptiveness": 0.25
}
```

### 8.7 Thảo Đỗ — Chủ quán tạp hóa

**Agent ID:** `agent_shop_owner`  
**Starting location:** `market`  
**Personality:** quen biết nhiều người, thực dụng, thích chuyện thời sự  
**Private motive:** Muốn giữ khách và trở thành điểm trao đổi tin của khu phố.

**Directly knows / witnessed**
- Nhìn thấy xe cứu hỏa và xe cấp cứu chạy qua; không thấy hiện trường.

**Does not know**
- Không biết thương vong hay nguyên nhân.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.58,
  "credibility": 0.62,
  "skepticism": 0.40,
  "sociality": 0.90,
  "deception_tendency": 0.12,
  "receptiveness": 0.68
}
```

### 8.8 Bác Dũng — Cư dân khu Riverside

**Agent ID:** `agent_resident`  
**Starting location:** `residential`  
**Personality:** cẩn trọng nhưng lo lắng cho khu dân cư  
**Private motive:** Sợ ô nhiễm ảnh hưởng gia đình.

**Directly knows / witnessed**
- Thấy cột khói xám, ngửi thấy mùi khét điện/nhựa; không ngửi thấy mùi hóa chất lạ.

**Does not know**
- Không có thông tin bên trong nhà máy.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.36,
  "credibility": 0.58,
  "skepticism": 0.50,
  "sociality": 0.62,
  "deception_tendency": 0.04,
  "receptiveness": 0.57
}
```

### 8.9 Khoa Bùi — Blogger địa phương

**Agent ID:** `agent_blogger`  
**Starting location:** `news_office`  
**Personality:** nhanh, tham engagement, tự tin quá mức  
**Private motive:** Muốn bài đăng đầu tiên đạt lượng xem lớn; sẵn sàng dùng tiêu đề mập mờ.

**Directly knows / witnessed**
- Nhận ảnh khói từ Nhiếp ảnh gia và tin nhắn 'nghe nói có nổ' từ Hùng.

**Does not know**
- Không có tài liệu y tế hay cứu hỏa.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.82,
  "credibility": 0.48,
  "skepticism": 0.33,
  "sociality": 0.80,
  "deception_tendency": 0.35,
  "receptiveness": 0.65
}
```

### 8.10 Vy Hoàng — Influencer khu vực

**Agent ID:** `agent_influencer`  
**Starting location:** `public_square`  
**Personality:** cảm tính, rất mạnh về lan truyền, thích phản ứng trực tiếp  
**Private motive:** Giữ hình ảnh 'người nói thay cộng đồng'; sợ bỏ lỡ chủ đề hot.

**Directly knows / witnessed**
- Không chứng kiến; chỉ xem bài của Khoa và bình luận của cư dân.

**Does not know**
- Mọi dữ kiện trực tiếp.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.96,
  "credibility": 0.57,
  "skepticism": 0.25,
  "sociality": 0.95,
  "deception_tendency": 0.22,
  "receptiveness": 0.73
}
```

### 8.11 Tuấn Đặng — Tài xế taxi

**Agent ID:** `agent_taxi_driver`  
**Starting location:** `public_square`  
**Personality:** giao tiếp rộng, nghe nhiều nguồn, nhớ chi tiết trung bình  
**Private motive:** Muốn kể chuyện cho khách nhưng không có lợi ích trực tiếp.

**Directly knows / witnessed**
- Chở một công nhân rời khu vực; người này nói 'chỉ bị khói, chưa nghe ai chết'.

**Does not know**
- Không biết người công nhân đó có đủ thông tin hay không.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.43,
  "credibility": 0.51,
  "skepticism": 0.38,
  "sociality": 0.84,
  "deception_tendency": 0.08,
  "receptiveness": 0.71
}
```

### 8.12 An Nguyễn — Sinh viên

**Agent ID:** `agent_student`  
**Starting location:** `cafe`  
**Personality:** tò mò, online nhiều, thích kiểm chứng nhưng dễ FOMO  
**Private motive:** Muốn tìm ra câu chuyện thật để đăng vào nhóm trường.

**Directly knows / witnessed**
- Không thấy hiện trường; có thể truy ra timestamp và phiên bản bài đăng khác nhau.

**Does not know**
- Không biết dữ kiện hiện trường nếu chưa hỏi.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.47,
  "credibility": 0.46,
  "skepticism": 0.57,
  "sociality": 0.86,
  "deception_tendency": 0.03,
  "receptiveness": 0.76
}
```

### 8.13 Ngọc Trịnh — Phóng viên địa phương

**Agent ID:** `agent_journalist`  
**Starting location:** `news_office`  
**Personality:** hoài nghi, có quy trình xác minh, khó bị thuyết phục bởi nguồn đơn  
**Private motive:** Muốn xuất bản chính xác, không muốn bị blogger vượt mặt.

**Directly knows / witnessed**
- Biết cách liên hệ cứu hỏa/bệnh viện; ban đầu chưa có tài liệu.

**Does not know**
- Chưa trực tiếp chứng kiến.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.78,
  "credibility": 0.88,
  "skepticism": 0.92,
  "sociality": 0.55,
  "deception_tendency": 0.02,
  "receptiveness": 0.24
}
```

### 8.14 Sơn Vương — Quản lý nhà máy

**Agent ID:** `agent_factory_manager`  
**Starting location:** `factory_office`  
**Personality:** phòng thủ, kiểm soát thông tin, chịu áp lực danh tiếng  
**Private motive:** Muốn giảm thiệt hại danh tiếng; che giấu việc bảo trì bị trễ nhưng không che giấu người chết vì thực tế không có ai chết.

**Directly knows / witnessed**
- Biết không có báo cáo tử vong; biết ticket bảo trì tủ điện đã quá hạn 12 ngày.

**Does not know**
- Không tận mắt thấy thời điểm chập điện.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.68,
  "credibility": 0.55,
  "skepticism": 0.70,
  "sociality": 0.42,
  "deception_tendency": 0.45,
  "receptiveness": 0.28
}
```

### 8.15 Đức Phan — Kỹ thuật viên bảo trì

**Agent ID:** `agent_technician`  
**Starting location:** `factory_office`  
**Personality:** kỹ tính, ngại truyền thông, có mặc cảm trách nhiệm  
**Private motive:** Sợ bị quy trách nhiệm vì ticket sửa điện bị trì hoãn.

**Directly knows / witnessed**
- Biết tủ DB-4 có lỗi quá nhiệt; đã tạo maintenance ticket; ảnh sau cháy cho thấy điểm hồ quang điện.

**Does not know**
- Không biết ai bắt đầu tin đồn.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.34,
  "credibility": 0.86,
  "skepticism": 0.84,
  "sociality": 0.28,
  "deception_tendency": 0.16,
  "receptiveness": 0.30
}
```

### 8.16 Hà Lương — Y tá

**Agent ID:** `agent_nurse`  
**Starting location:** `hospital`  
**Personality:** chắc chắn, giữ bí mật bệnh nhân, thân thiện  
**Private motive:** Không tiết lộ danh tính bệnh nhân nhưng có thể xác nhận tổng số ca.

**Directly knows / witnessed**
- Biết không có tử vong và không có ca bỏng hóa chất.

**Does not know**
- Không biết nguồn cháy.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.39,
  "credibility": 0.84,
  "skepticism": 0.78,
  "sociality": 0.48,
  "deception_tendency": 0.01,
  "receptiveness": 0.31
}
```

### 8.17 Cô Hương — Người bán hàng rong

**Agent ID:** `agent_vendor`  
**Starting location:** `factory_gate`  
**Personality:** nhiều mối quan hệ, nhớ người tốt nhưng nhớ nội dung không hoàn hảo  
**Private motive:** Không có động cơ xấu; thích kể chuyện.

**Directly knows / witnessed**
- Thấy xe cấp cứu đưa vài người đi nhưng tất cả đều tự đi hoặc ngồi được.

**Does not know**
- Không biết tình trạng sau đó.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.49,
  "credibility": 0.50,
  "skepticism": 0.32,
  "sociality": 0.93,
  "deception_tendency": 0.06,
  "receptiveness": 0.78
}
```

### 8.18 Thầy Bình — Giáo viên nghỉ hưu

**Agent ID:** `agent_retired_teacher`  
**Starting location:** `residential`  
**Personality:** kiên nhẫn, kiểm chứng trước khi tin, có tiếng nói trong cộng đồng  
**Private motive:** Muốn tránh hoảng loạn trong khu dân cư.

**Directly knows / witnessed**
- Không có quan sát trực tiếp đáng kể; có năng lực đánh giá mâu thuẫn giữa các lời kể.

**Does not know**
- Thiếu dữ kiện gốc.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.61,
  "credibility": 0.79,
  "skepticism": 0.82,
  "sociality": 0.50,
  "deception_tendency": 0.00,
  "receptiveness": 0.29
}
```

### 8.19 Phúc Trần — Tài xế giao hàng

**Agent ID:** `agent_delivery_driver`  
**Starting location:** `market`  
**Personality:** vội, chú ý đường đi, dễ lặp lại tin nghe được  
**Private motive:** Chỉ muốn hoàn thành chuyến hàng.

**Directly knows / witnessed**
- Đi ngang nhà máy trước cháy 10 phút và thấy xe kỹ thuật đỗ gần kho; không biết lý do.

**Does not know**
- Không thấy thời điểm cháy.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.41,
  "credibility": 0.49,
  "skepticism": 0.36,
  "sociality": 0.72,
  "deception_tendency": 0.05,
  "receptiveness": 0.73
}
```

### 8.20 Cô Liên — Trưởng nhóm cộng đồng

**Agent ID:** `agent_community_leader`  
**Starting location:** `public_square`  
**Personality:** trách nhiệm, ưu tiên ổn định, có mạng lưới rộng  
**Private motive:** Muốn đưa thông tin hướng dẫn cộng đồng đáng tin.

**Directly knows / witnessed**
- Có kênh liên hệ chính quyền/cứu hỏa nhưng phản hồi đến chậm.

**Does not know**
- Không biết nội bộ nhà máy.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.88,
  "credibility": 0.82,
  "skepticism": 0.69,
  "sociality": 0.82,
  "deception_tendency": 0.03,
  "receptiveness": 0.44
}
```

### 8.21 Tùng Lâm — Nhiếp ảnh gia tự do

**Agent ID:** `agent_photographer`  
**Starting location:** `factory_gate`  
**Personality:** quan sát hình ảnh tốt, thích khoảnh khắc gây ấn tượng  
**Private motive:** Muốn bán ảnh; crop ảnh có thể khiến khói trông nghiêm trọng hơn.

**Directly knows / witnessed**
- Chụp chuỗi ảnh timestamp: khói bắt đầu ở khu kho, không có fireball; ảnh cho thấy cửa kho còn nguyên kết cấu.

**Does not know**
- Không biết tình trạng người bị thương.

**Dialogue direction**
- Never speak with knowledge outside direct observations, memories received, public posts, or evidence explicitly revealed during play.
- If uncertain, use uncertainty language.
- If asked a leading question, do not automatically agree.
- If the character has a deception motive, they may evade or frame facts, but may not conjure server-hidden facts.

**Initial numeric profile**

```json
{
  "influence": 0.66,
  "credibility": 0.63,
  "skepticism": 0.54,
  "sociality": 0.67,
  "deception_tendency": 0.10,
  "receptiveness": 0.58
}
```

# 9. CLAIM CATALOG
| Claim ID | Text | Server classification | Role |
|---|---|---|---|
| `claim_fire_occurred` | Có cháy tại kho Riverside | `true` | `fact` |
| `claim_electrical_origin` | Điểm khởi phát nằm gần tủ phân phối điện DB-4 | `true` | `fact` |
| `claim_electrical_fault_preexisting` | DB-4 đã có cảnh báo quá nhiệt trước sự cố | `true` | `fact` |
| `claim_minor_injuries` | Có một số ca bị thương nhẹ do khói/sơ tán | `true` | `fact` |
| `claim_zero_fatalities` | Không có người tử vong | `true` | `fact` |
| `claim_no_major_explosion` | Không có vụ nổ lớn | `true` | `fact` |
| `claim_no_chemical_release` | Không có bằng chứng rò rỉ hóa chất nguy hiểm | `true` | `fact` |
| `claim_chemical_explosion` | Nhà máy xảy ra nổ hóa chất | `false` | `primary_false` |
| `claim_multiple_deaths` | Nhiều công nhân đã chết | `false` | `primary_false` |
| `claim_company_coverup` | Công ty đang che giấu số người chết | `false` | `primary_false` |
| `claim_toxic_cloud` | Đám mây độc đang lan vào khu dân cư | `false` | `secondary_false` |
| `claim_blogger_had_verified_deaths` | Blogger có nguồn xác nhận người chết | `false` | `secondary_false` |
| `claim_rumor_escalated_without_new_source` | Tin đồn tăng mức nghiêm trọng dù không có nguồn mới | `true` | `meta` |

## 9.1 Security boundary

`classification` is server-only scenario metadata.

Never include it in:
- normal REST responses,
- WebSocket state,
- frontend stores,
- LLM context,
- normal logs exposed to the browser.

It may appear only in:
- scenario validation,
- server-side scoring,
- tests,
- explicit development diagnostic endpoints guarded by development mode.

# 10. RELATIONSHIP GRAPH
Relationships are **directed**. If only one direction is listed in scenario data, the loader may create a default reverse edge at a lower or configured value, or require both edges explicitly. For V1, prefer explicit directed rows.
| From | To | Trust | Reason |
|---|---|---:|---|
| `agent_factory_worker_a` | `agent_factory_worker_b` | 0.78 | Đồng nghiệp thân |
| `agent_factory_worker_b` | `agent_shop_owner` | 0.72 | Thường uống nước và kể chuyện |
| `agent_shop_owner` | `agent_vendor` | 0.81 | Quan hệ hàng xóm |
| `agent_blogger` | `agent_influencer` | 0.74 | Từng hợp tác nội dung |
| `agent_student` | `agent_blogger` | 0.41 | Theo dõi nhưng không hoàn toàn tin |
| `agent_journalist` | `agent_blogger` | 0.26 | Cạnh tranh nguồn tin |
| `agent_community_leader` | `agent_retired_teacher` | 0.83 | Tin tưởng lẫn nhau |
| `agent_community_leader` | `agent_journalist` | 0.76 | Nguồn thông tin đáng tin |
| `agent_factory_manager` | `agent_technician` | 0.52 | Quan hệ công việc căng thẳng |
| `agent_firefighter` | `agent_paramedic` | 0.88 | Phối hợp nghề nghiệp |
| `agent_paramedic` | `agent_nurse` | 0.92 | Đồng nghiệp |
| `agent_resident` | `agent_retired_teacher` | 0.73 | Hàng xóm |
| `agent_photographer` | `agent_blogger` | 0.66 | Bán ảnh/tin |
| `agent_taxi_driver` | `agent_vendor` | 0.63 | Khách quen |
| `agent_delivery_driver` | `agent_shop_owner` | 0.69 | Giao hàng thường xuyên |

## 10.1 Relationship use in propagation

Trust is not a universal “friendship” score.

For information propagation it represents:

> “How much does listener B currently trust information from speaker A?”

Therefore:

- A close friend may have high trust.
- A journalist may assign low trust to a sensational blogger.
- A nurse may trust a paramedic professionally.
- Trust can change after verified corrections or proven false statements.

Default missing trust:

```text
default_relationship_trust = 0.50
```

Never infer trust from occupation at runtime. Occupation may only influence scenario seed values.

# 11. EVIDENCE BIBLE
| Evidence ID | Name | Location | Reliability | Supports | Contradicts | Unlock |
|---|---|---|---:|---|---|---|
| `ev_cctv_gate` | CCTV cổng nhà máy | `factory_security_room` | 0.94 | claim_fire_occurred, claim_no_major_explosion | claim_chemical_explosion, claim_mass_casualties | Nói chuyện với Minh hoặc inspect phòng bảo vệ sau khi có quyền truy cập. |
| `ev_fire_report` | Biên bản sơ bộ đội cứu hỏa | `fire_station` | 0.98 | claim_electrical_origin, claim_fire_occurred | claim_chemical_explosion | Phỏng vấn Quang với credibility player >= 0.45 hoặc đã có ảnh hiện trường. |
| `ev_hospital_summary` | Tóm tắt tiếp nhận bệnh viện | `hospital` | 0.99 | claim_zero_fatalities, claim_minor_injuries | claim_multiple_deaths, claim_chemical_exposure | Phỏng vấn Mai/Hà; không lộ thông tin cá nhân. |
| `ev_maintenance_ticket` | Ticket bảo trì DB-4 | `factory_office` | 0.91 | claim_electrical_fault_preexisting | claim_random_chemical_explosion | Tìm trong văn phòng kỹ thuật sau khi Đức tiết lộ mã DB-4 hoặc Sơn cho phép. |
| `ev_panel_photo` | Ảnh tủ điện sau cháy | `factory_storage` | 0.93 | claim_electrical_origin | claim_chemical_explosion | Inspect kho sau khi Fire Station đánh dấu khu vực an toàn. |
| `ev_photo_sequence` | Chuỗi ảnh timestamp của Tùng | `factory_gate` | 0.88 | claim_fire_occurred, claim_no_major_explosion | claim_large_explosion | Phỏng vấn Tùng hoặc thuyết phục mua/nhận ảnh. |
| `ev_worker_testimony` | Lời khai trực tiếp của Lan | `factory_yard` | 0.78 | claim_electrical_origin, claim_zero_observed_deaths | claim_large_explosion | Phỏng vấn Lan khi cô đã bình tĩnh; hỏi cụ thể 'cô nhìn thấy gì' thay vì 'có nổ không'. |
| `ev_dispatch_log` | Nhật ký điều xe cấp cứu | `hospital` | 0.96 | claim_minor_injuries, claim_zero_fatalities | claim_mass_casualties | Có hospital summary rồi và hỏi Mai về quy mô điều xe. |
| `ev_blog_versions` | Lịch sử chỉnh sửa bài blog | `news_office` | 0.86 | claim_rumor_escalated_without_new_source | claim_blogger_had_verified_deaths | Dùng terminal ở News Office sau khi An chỉ ra cached versions. |
| `ev_air_sensor` | Cảm biến môi trường khu dân cư | `residential` | 0.84 | claim_no_chemical_release | claim_toxic_cloud | Được Thầy Bình chỉ vị trí sau khi hỏi về mùi/không khí. |

## 11.1 Detailed evidence descriptions

### Evidence 1: CCTV cổng nhà máy

**ID:** `ev_cctv_gate`  
**Location:** `factory_security_room`  
**Reliability:** `0.94`

**What the player obtains**

Video 17:56–18:05 cho thấy không có sóng nổ/fireball; người sơ tán tự chạy ra.

**Supports**
`claim_fire_occurred, claim_no_major_explosion`

**Contradicts**
`claim_chemical_explosion, claim_mass_casualties`

**Discovery condition**
Nói chuyện với Minh hoặc inspect phòng bảo vệ sau khi có quyền truy cập.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 2: Biên bản sơ bộ đội cứu hỏa

**ID:** `ev_fire_report`  
**Location:** `fire_station`  
**Reliability:** `0.98`

**What the player obtains**

Ghi điểm cháy tập trung quanh tủ điện DB-4; không phát hiện dấu vết nổ hóa chất.

**Supports**
`claim_electrical_origin, claim_fire_occurred`

**Contradicts**
`claim_chemical_explosion`

**Discovery condition**
Phỏng vấn Quang với credibility player >= 0.45 hoặc đã có ảnh hiện trường.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 3: Tóm tắt tiếp nhận bệnh viện

**ID:** `ev_hospital_summary`  
**Location:** `hospital`  
**Reliability:** `0.99`

**What the player obtains**

4 ca nhẹ, 0 tử vong, 0 bỏng hóa chất.

**Supports**
`claim_zero_fatalities, claim_minor_injuries`

**Contradicts**
`claim_multiple_deaths, claim_chemical_exposure`

**Discovery condition**
Phỏng vấn Mai/Hà; không lộ thông tin cá nhân.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 4: Ticket bảo trì DB-4

**ID:** `ev_maintenance_ticket`  
**Location:** `factory_office`  
**Reliability:** `0.91`

**What the player obtains**

Ticket cảnh báo quá nhiệt và đánh lửa bất thường, quá hạn 12 ngày.

**Supports**
`claim_electrical_fault_preexisting`

**Contradicts**
`claim_random_chemical_explosion`

**Discovery condition**
Tìm trong văn phòng kỹ thuật sau khi Đức tiết lộ mã DB-4 hoặc Sơn cho phép.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 5: Ảnh tủ điện sau cháy

**ID:** `ev_panel_photo`  
**Location:** `factory_storage`  
**Reliability:** `0.93`

**What the player obtains**

Ảnh cho thấy dấu hồ quang điện và cháy lan cục bộ từ tủ phân phối.

**Supports**
`claim_electrical_origin`

**Contradicts**
`claim_chemical_explosion`

**Discovery condition**
Inspect kho sau khi Fire Station đánh dấu khu vực an toàn.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 6: Chuỗi ảnh timestamp của Tùng

**ID:** `ev_photo_sequence`  
**Location:** `factory_gate`  
**Reliability:** `0.88`

**What the player obtains**

6 ảnh liên tiếp chứng minh khói tăng dần, không có blast wave hay fireball.

**Supports**
`claim_fire_occurred, claim_no_major_explosion`

**Contradicts**
`claim_large_explosion`

**Discovery condition**
Phỏng vấn Tùng hoặc thuyết phục mua/nhận ảnh.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 7: Lời khai trực tiếp của Lan

**ID:** `ev_worker_testimony`  
**Location:** `factory_yard`  
**Reliability:** `0.78`

**What the player obtains**

Lan nghe tiếng tách/nổ nhỏ tại tủ điện và nhìn thấy tia lửa; cô không thấy người chết.

**Supports**
`claim_electrical_origin, claim_zero_observed_deaths`

**Contradicts**
`claim_large_explosion`

**Discovery condition**
Phỏng vấn Lan khi cô đã bình tĩnh; hỏi cụ thể 'cô nhìn thấy gì' thay vì 'có nổ không'.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 8: Nhật ký điều xe cấp cứu

**ID:** `ev_dispatch_log`  
**Location:** `hospital`  
**Reliability:** `0.96`

**What the player obtains**

Hai xe được điều động, chở tổng cộng bốn ca không nguy kịch.

**Supports**
`claim_minor_injuries, claim_zero_fatalities`

**Contradicts**
`claim_mass_casualties`

**Discovery condition**
Có hospital summary rồi và hỏi Mai về quy mô điều xe.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 9: Lịch sử chỉnh sửa bài blog

**ID:** `ev_blog_versions`  
**Location:** `news_office`  
**Reliability:** `0.86`

**What the player obtains**

Bản đầu: 'nghe tiếng nổ'; bản sau: 'nghi nổ hóa chất'; bản viral: 'có người chết' dù không thêm nguồn mới.

**Supports**
`claim_rumor_escalated_without_new_source`

**Contradicts**
`claim_blogger_had_verified_deaths`

**Discovery condition**
Dùng terminal ở News Office sau khi An chỉ ra cached versions.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.

### Evidence 10: Cảm biến môi trường khu dân cư

**ID:** `ev_air_sensor`  
**Location:** `residential`  
**Reliability:** `0.84`

**What the player obtains**

Không ghi nhận spike hóa chất nguy hiểm tương ứng thời điểm cháy; chỉ tăng PM do khói.

**Supports**
`claim_no_chemical_release`

**Contradicts**
`claim_toxic_cloud`

**Discovery condition**
Được Thầy Bình chỉ vị trí sau khi hỏi về mùi/không khí.

**Design note**
The UI may summarize what the evidence means, but the server stores explicit claim links. Evidence acquisition must never silently reveal the server's `true/false` classification labels.


# 12. PLAYER ACTION ECONOMY

Actions cost **game time**, creating opportunity cost.

Initial tuning:

| Action | Game-time cost | Notes |
|---|---:|---|
| Select/read basic NPC card | 0 | UI only |
| Travel to adjacent location | 20–45 sec | Depends on waypoint distance |
| Interview short exchange | 30 sec | Charged after NPC answer |
| Follow-up interview question | 20 sec | Same conversation |
| Inspect location | 30–60 sec | Determined by scenario |
| Acquire evidence | 10–30 sec | Some evidence requires prior unlock |
| Review notebook | 0 | Pauses only if accessibility mode enables pause |
| Publish correction | 45 sec | Includes composing/verifying |
| Submit Ground Truth | 15 sec | Terminal if correct/qualified |

Do not blindly stop world simulation when the player opens a panel.

Recommended default:
- normal panels: simulation continues,
- accessibility “pause while reading” option: allowed, but scoring should mark the run as assisted.

---

# 13. BELIEF MODEL

## 13.1 Agent belief row

For every relevant `(agent_id, claim_id)` pair:

```go
type AgentBelief struct {
    AgentID              string
    ClaimID              string
    Confidence           float64 // [0,1]
    FirstSourceAgentID   *string
    LastSourceAgentID    *string
    EvidenceResistance   float64
    UpdatedAtGameSecond  int64
    Revision             int64
}
```

## 13.2 Initial beliefs

Do not initialize every claim to `0.5`.

Unknown is not the same as “50% convinced”.

Use a baseline such as:

```text
unknown belief confidence = 0.10
```

Then direct observations seed confidence based on observation reliability.

Example:

Lan directly sees sparks near DB-4:

```text
claim_electrical_origin = 0.72
```

Minh hears a sharp bang but does not see the panel:

```text
claim_large_explosion = 0.10
claim_fire_occurred   = 0.90
```

## 13.3 Propagation strength

Baseline:

```text
speaker_confidence_weight   = 0.30
relationship_trust_weight   = 0.25
speaker_influence_weight    = 0.15
speaker_credibility_weight  = 0.15
listener_receptiveness      = 0.15
```

Formula:

```text
base =
  speaker_confidence * 0.30 +
  relationship_trust * 0.25 +
  speaker_influence * 0.15 +
  speaker_credibility * 0.15 +
  listener_receptiveness * 0.15
```

Weights must sum to 1.0 and startup validation must reject a material mismatch.

## 13.4 Skepticism modifier

```text
skepticism_multiplier = 1 - (listener_skepticism * skepticism_effect)
```

Suggested:

```text
skepticism_effect = 0.55
```

## 13.5 Direct contradictory observation

If listener has a high-reliability direct observation contradicting a claim:

```text
direct_evidence_multiplier = 0.35
```

Do not make it zero. Human characters can still become confused.

## 13.6 Repeated-source diminishing return

Receiving the same claim from the same source family repeatedly should not count like independent confirmation.

Maintain provenance graph or at minimum:

```text
origin_memory_id
root_source_agent_id
```

If the new message has the same root source as a previous memory:

```text
same_root_source_multiplier = 0.45
```

This is essential. Without it, one influencer echoed by ten people creates fake “independent confirmation” mechanically.

## 13.7 Update formula

For a positive assertion:

```text
delta =
  (1 - old_confidence)
  * propagation_strength
  * belief_update_rate
  * skepticism_multiplier
  * evidence_modifier
  * provenance_modifier

new_confidence = clamp(old_confidence + delta, 0, 1)
```

Suggested:

```text
belief_update_rate = 0.38
```

For an explicit contradiction/correction, use a separate correction rule rather than negative propagation through the exact same formula.

## 13.8 Determinism

All non-LLM simulation randomness must use a session seed.

Never call a global random source directly inside domain rules.

Store:

```text
simulation_seed
simulation_revision
```

A test run with:
- same scenario,
- same seed,
- same player actions,

must produce the same belief transitions when LLM is disabled.

---

# 14. PUBLIC CORRECTION MODEL

The player can publish only a claim-specific correction.

Input:

```json
{
  "challenged_claim_id": "claim_multiple_deaths",
  "evidence_ids": ["ev_hospital_summary", "ev_dispatch_log"],
  "message": "Hospital records show four minor cases and zero fatalities."
}
```

The text is not the authority. The selected evidence is.

## 14.1 Correction strength

Baseline:

```text
evidence_quality      0.45
source_diversity      0.20
player_credibility    0.20
corroboration         0.15
```

A correction cannot use undiscovered evidence.

## 14.2 Distribution

Publishing does not instantly hit all NPCs.

Model one of:

- public square announcement,
- local social post,
- journalist amplification,
- community-leader amplification.

Each has reach.

V1 simplest model:

```text
base_public_reach = 0.55
```

Then agents with high social connectivity may relay it.

## 14.3 Backfire

Do not implement a complex ideological backfire mechanic in V1.

A weak correction simply has low effect and can reduce player credibility if it contains an incorrect assertion.

## 14.4 Player credibility

Start:

```text
player_credibility = 0.50
```

Credibility increases when:
- published correction is supported by strong evidence,
- later official evidence corroborates it.

Credibility decreases when:
- player publishes a claim contradicted by strong discovered evidence,
- player claims certainty without minimum evidence.

Clamp `[0,1]`.

---

# 15. EVIDENCE STRENGTH

Evidence Strength is not “number of clues”.

```text
evidence_strength =
  coverage_score    * 0.50 +
  reliability_score * 0.30 +
  diversity_score   * 0.20
```

## 15.1 Coverage

Required Ground Truth fields:

```text
event_type
cause_category
major_explosion
fatalities
```

Each field is covered when at least one qualifying discovered evidence item supports the correct answer.

Strong independent corroboration may increase field confidence but coverage remains capped.

## 15.2 Reliability

Weighted mean of evidence reliability, but prevent evidence spam.

Only relevant evidence for submitted fields contributes.

## 15.3 Diversity

Evidence source categories:

```text
direct_witness
official_response
medical
technical_record
visual_record
digital_trace
environmental_sensor
```

More independent categories increase diversity.

---

# 16. DIALOGUE AND LLM RULES

## 16.1 The LLM's job

The LLM controls:

- wording,
- tone,
- personality expression,
- conversational coherence,
- whether an allowed deceptive character evades or lies within its knowledge model.

The LLM does **not** control:

- what happened,
- which claims exist,
- private server classification,
- exact belief mutations,
- win/loss,
- evidence discovery unless a preconfigured rule allows the NPC to reveal an evidence ID.

## 16.2 Context builder output

Never dump an Agent database row directly into a prompt.

Build a purpose-specific context:

```json
{
  "identity": {
    "name": "Minh Trần",
    "role": "Security Guard"
  },
  "traits": {
    "skepticism": 0.72,
    "deception_tendency": 0.10
  },
  "current_location": "factory_gate",
  "known_claims": [
    {
      "claim_id": "claim_fire_occurred",
      "confidence": 0.91,
      "basis": "direct_observation"
    }
  ],
  "relevant_memories": [],
  "relationship_to_player": 0.50,
  "conversation_history": [],
  "allowed_reveal_evidence_ids": ["ev_cctv_gate"]
}
```

Do not include classifications.

## 16.3 Required model output

```json
{
  "intent": "answer",
  "utterance": "Tôi thấy khói từ khu kho. Có một tiếng nổ nhỏ, nhưng tôi không thấy kiểu vụ nổ lớn nào cả.",
  "referenced_claim_ids": [
    "claim_fire_occurred",
    "claim_no_major_explosion"
  ],
  "revealed_evidence_ids": [],
  "emotion": "calm",
  "certainty": 0.78
}
```

Allowed intents:

```text
answer
deflect
lie
refuse
ask_question
```

Allowed emotions:

```text
neutral
calm
uncertain
afraid
angry
excited
defensive
sad
```

## 16.4 Knowledge guard

Reject a model response if:
- it references an unknown claim,
- it reveals an evidence ID not in `allowed_reveal_evidence_ids`,
- it invents a person/location identifier,
- it includes system/debug text,
- its JSON fails validation.

Retry at most configured count.

After retry exhaustion, use a deterministic template response.

## 16.5 Dialogue latency

Simulation must never stop waiting for dialogue.

The interview UI may show:
> “Minh đang suy nghĩ…”

while the simulation continues.

Use a bounded worker queue.

Defaults:

```text
llm_workers              = 3
llm_queue_capacity       = 50
llm_request_timeout      = 15s
llm_retry_count          = 1
llm_max_output_chars     = 1200
```

---

# 17. NPC AUTONOMY SCHEDULER

## 17.1 Scheduler responsibilities

The scheduler decides:
- who wants to talk,
- which nearby/socially connected target is selected,
- which known claim is salient,
- whether the conversation is skipped due to cooldown,
- when a public post is emitted.

It does not generate dialogue wording.

## 17.2 Salience score

A claim becomes more likely to be shared when:
- confidence is high,
- emotional relevance is high,
- novelty is high,
- the NPC is social,
- the claim is currently trending.

Example:

```text
salience =
  confidence        * 0.35 +
  novelty           * 0.20 +
  emotional_weight  * 0.15 +
  sociality         * 0.15 +
  trending_score    * 0.15
```

## 17.3 Conversation candidate constraints

Do not schedule conversation when:
- either NPC is already busy,
- distance/social relation threshold fails,
- same pair cooldown is active,
- session terminal state reached,
- per-tick conversation budget exceeded.

Defaults:

```text
max_new_npc_conversations_per_tick = 2
pair_cooldown_game_seconds          = 90
agent_share_cooldown_game_seconds   = 40
```

## 17.4 Why not LLM-plan every NPC

The project is evaluating multi-agent social dynamics, not how many API calls can be made.

Deterministic scheduling provides:
- reproducible tests,
- lower cost,
- explicit game design,
- fast simulation,
- observable engineering behavior.

---

# 18. SESSION STATE MACHINE

```text
CREATED
  ↓
INITIALIZING
  ↓
RUNNING
  ├──→ WON
  ├──→ LOST_FALSE_BELIEF
  ├──→ LOST_TIMEOUT (optional)
  └──→ ABORTED
```

Rules:
- only `RUNNING` sessions tick,
- terminal sessions reject state-mutating gameplay actions,
- read endpoints still work after terminal state,
- refresh must reload persisted terminal result,
- WebSocket reconnect receives latest snapshot plus sequence cursor.

---

# 19. REAL-TIME EVENT ORDER

Each session has a monotonically increasing `sequence`.

Example:

```json
{
  "sequence": 184,
  "type": "belief_stats.updated",
  "session_id": "uuid",
  "game_second": 1320,
  "occurred_at": "2026-09-10T02:30:00Z",
  "payload": {
    "primary_false_narrative_ratio": 0.47
  }
}
```

Client rules:
1. ignore duplicate `sequence`,
2. if sequence gap > 1, request/resync snapshot,
3. never derive authoritative belief ratio locally.

---

# 20. BACKEND TECHNOLOGY DECISIONS

Use:
- Go 1.27.x.
- `net/http` server with `go-chi/chi/v5` for composable routing/middleware.
- `github.com/coder/websocket` for WebSocket transport.
- PostgreSQL.
- `github.com/jackc/pgx/v5` / `pgxpool`.
- `sqlc` for type-safe query generation.
- versioned SQL migrations.
- `log/slog` for structured logs.
- `context.Context` everywhere at I/O boundaries.
- standard `testing`, `httptest`; `testcontainers-go` only for DB integration tests if needed.

Do not add an ORM.

Do not use Gin/Fiber unless the owner explicitly changes the decision.

Rationale:
- the project benefits from explicit domain/application boundaries,
- HTTP complexity is modest,
- pgx + sqlc keeps SQL visible and testable,
- WebSocket is a narrow transport concern,
- Go standard library should remain prominent.

---

# 21. GO CODE STYLE

## 21.1 Module

Recommended module path:

```text
github.com/<owner>/city-of-lies/backend
```

Replace `<owner>` only when the repository owner is known.

## 21.2 Package naming

Use short domain names:

```text
domain
belief
evidence
outcome
simulation
scenario
postgres
realtime
handlers
```

Do not create packages named:

```text
utils
helpers
common
misc
base
```

If something has no clear domain, reconsider its ownership.

## 21.3 Interfaces

Define interfaces at the consumer boundary.

Bad:

```go
type GenericRepository[T any] interface { ... }
```

Preferred:

```go
type SessionRepository interface {
    Create(ctx context.Context, session domain.GameSession) error
    Get(ctx context.Context, id string) (domain.GameSession, error)
}
```

## 21.4 Errors

Wrap errors with context:

```go
return fmt.Errorf("load scenario %s: %w", id, err)
```

Use domain sentinel/type errors only for cases application code branches on.

Do not expose raw SQL error strings to clients.

## 21.5 Concurrency

Every goroutine must have:
- an owner,
- cancellation path,
- bounded queue where applicable,
- shutdown behavior.

No anonymous “fire and forget” goroutines in handlers.

## 21.6 Race safety

CI should include:

```bash
go test -race ./...
```

at least on a dedicated workflow/job if runtime cost is acceptable.

---

# 22. CONFIGURATION MODEL

Configuration hierarchy:

1. code defaults,
2. environment variables,
3. command-line override only if explicitly needed for developer tools.

No secret config file committed.

Strong config structs:

```go
type Config struct {
    HTTP       HTTPConfig
    Database   DatabaseConfig
    Simulation SimulationConfig
    Game       GameConfig
    LLM        LLMConfig
    CORS       CORSConfig
}
```

Validation runs before server start.

Invalid examples that must fail fast:
- threshold outside `[0,1]`,
- negative tick duration,
- zero WebSocket queue capacity,
- LLM enabled without base URL/model/key when provider requires them,
- DB DSN missing,
- CORS wildcard in production when credentials enabled.


# 23. REPOSITORY STRUCTURE

```text
city-of-lies/
├── README.md
├── LICENSE
├── .gitignore
├── .editorconfig
├── .env.example
├── Makefile
├── docker-compose.yml
├── docker-compose.dev.yml
│
├── docs/
│   ├── GAME_BIBLE.md
│   ├── SCENARIO_RIVERSIDE.md
│   ├── ARCHITECTURE.md
│   ├── PROTOCOL.md
│   ├── DATABASE.md
│   ├── TEST_PLAN.md
│   ├── IMPLEMENTATION_STATUS.md
│   └── ADR/
│       ├── 0001-go-modular-monolith.md
│       ├── 0002-deterministic-simulation.md
│       ├── 0003-llm-not-source-of-truth.md
│       ├── 0004-websocket-over-polling.md
│       └── 0005-pgx-sqlc-over-orm.md
│
├── data/
│   └── scenarios/
│       └── riverside-factory/
│           ├── scenario.json
│           ├── locations.json
│           ├── agents.json
│           ├── claims.json
│           ├── observations.json
│           ├── relationships.json
│           ├── evidence.json
│           ├── public-events.json
│           └── truth-form.json
│
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── sqlc.yaml
│   ├── Dockerfile
│   ├── .golangci.yml
│   │
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go
│   │   ├── sim/
│   │   │   └── main.go
│   │   └── scenario-check/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   ├── defaults.go
│   │   │   └── config_test.go
│   │   │
│   │   ├── domain/
│   │   │   ├── ids.go
│   │   │   ├── session.go
│   │   │   ├── agent.go
│   │   │   ├── claim.go
│   │   │   ├── belief.go
│   │   │   ├── memory.go
│   │   │   ├── relationship.go
│   │   │   ├── evidence.go
│   │   │   ├── location.go
│   │   │   ├── player_investigation.go
│   │   │   ├── world_event.go
│   │   │   ├── dialogue.go
│   │   │   └── errors.go
│   │   │
│   │   ├── game/
│   │   │   ├── belief/
│   │   │   │   ├── propagation.go
│   │   │   │   ├── propagation_test.go
│   │   │   │   ├── narrative.go
│   │   │   │   ├── narrative_test.go
│   │   │   │   ├── correction.go
│   │   │   │   └── correction_test.go
│   │   │   ├── evidence/
│   │   │   │   ├── strength.go
│   │   │   │   └── strength_test.go
│   │   │   ├── outcome/
│   │   │   │   ├── evaluator.go
│   │   │   │   └── evaluator_test.go
│   │   │   ├── clock/
│   │   │   │   ├── clock.go
│   │   │   │   └── fake_clock.go
│   │   │   ├── scheduler/
│   │   │   │   ├── scheduler.go
│   │   │   │   ├── policy.go
│   │   │   │   ├── random.go
│   │   │   │   └── scheduler_test.go
│   │   │   ├── movement/
│   │   │   │   ├── waypoints.go
│   │   │   │   └── movement.go
│   │   │   ├── dialogue/
│   │   │   │   ├── context_builder.go
│   │   │   │   ├── guard.go
│   │   │   │   ├── queue.go
│   │   │   │   ├── template_provider.go
│   │   │   │   ├── openai_compatible.go
│   │   │   │   └── dialogue_test.go
│   │   │   ├── scenario/
│   │   │   │   ├── loader.go
│   │   │   │   ├── validator.go
│   │   │   │   ├── model.go
│   │   │   │   └── validator_test.go
│   │   │   └── simulation/
│   │   │       ├── engine.go
│   │   │       ├── tick.go
│   │   │       ├── manager.go
│   │   │       ├── events.go
│   │   │       ├── snapshot.go
│   │   │       └── simulation_test.go
│   │   │
│   │   ├── application/
│   │   │   ├── ports/
│   │   │   │   ├── repositories.go
│   │   │   │   ├── realtime.go
│   │   │   │   └── dialogue.go
│   │   │   ├── dto/
│   │   │   │   ├── session.go
│   │   │   │   ├── agent.go
│   │   │   │   ├── notebook.go
│   │   │   │   └── requests.go
│   │   │   └── services/
│   │   │       ├── start_session.go
│   │   │       ├── get_state.go
│   │   │       ├── interview.go
│   │   │       ├── inspect.go
│   │   │       ├── notebook.go
│   │   │       ├── publish_correction.go
│   │   │       └── submit_truth.go
│   │   │
│   │   ├── infrastructure/
│   │   │   ├── postgres/
│   │   │   │   ├── pool.go
│   │   │   │   ├── tx.go
│   │   │   │   ├── repository/
│   │   │   │   │   ├── sessions.go
│   │   │   │   │   ├── agents.go
│   │   │   │   │   ├── beliefs.go
│   │   │   │   │   ├── memories.go
│   │   │   │   │   ├── evidence.go
│   │   │   │   │   └── events.go
│   │   │   │   └── sqlc/
│   │   │   │       └── [generated files]
│   │   │   └── realtime/
│   │   │       ├── hub.go
│   │   │       ├── client.go
│   │   │       ├── messages.go
│   │   │       └── publisher.go
│   │   │
│   │   └── http/
│   │       ├── router.go
│   │       ├── middleware/
│   │       │   ├── request_id.go
│   │       │   ├── logging.go
│   │       │   ├── recovery.go
│   │       │   └── cors.go
│   │       ├── handlers/
│   │       │   ├── health.go
│   │       │   ├── sessions.go
│   │       │   ├── agents.go
│   │       │   ├── investigation.go
│   │       │   └── ws.go
│   │       └── respond/
│   │           └── json.go
│   │
│   └── db/
│       ├── migrations/
│       │   ├── 000001_init.up.sql
│       │   ├── 000001_init.down.sql
│       │   ├── 000002_indexes.up.sql
│       │   └── 000002_indexes.down.sql
│       └── queries/
│           ├── sessions.sql
│           ├── agents.sql
│           ├── beliefs.sql
│           ├── memories.sql
│           ├── evidence.sql
│           ├── conversations.sql
│           └── events.sql
│
├── frontend/
│   ├── package.json
│   ├── package-lock.json
│   ├── next.config.ts
│   ├── tsconfig.json
│   ├── eslint.config.mjs
│   ├── Dockerfile
│   ├── .env.local.example
│   ├── public/
│   │   ├── models/
│   │   ├── textures/
│   │   ├── icons/
│   │   └── audio/
│   └── src/
│       ├── app/
│       │   ├── globals.css
│       │   ├── layout.tsx
│       │   ├── page.tsx
│       │   └── game/
│       │       └── [sessionId]/
│       │           └── page.tsx
│       ├── components/
│       │   ├── game/
│       │   │   └── GameShell.tsx
│       │   ├── world/
│       │   │   ├── GameCanvas.tsx
│       │   │   ├── CityScene.tsx
│       │   │   ├── CameraRig.tsx
│       │   │   ├── District.tsx
│       │   │   ├── LocationBuilding.tsx
│       │   │   ├── AgentAvatar.tsx
│       │   │   ├── AgentLayer.tsx
│       │   │   ├── RumorPulse.tsx
│       │   │   └── WorldInteraction.tsx
│       │   ├── hud/
│       │   │   ├── GameHud.tsx
│       │   │   ├── BeliefMeter.tsx
│       │   │   ├── EvidenceMeter.tsx
│       │   │   ├── CredibilityMeter.tsx
│       │   │   ├── EventFeed.tsx
│       │   │   └── SelectedAgentPanel.tsx
│       │   ├── interview/
│       │   │   ├── InterviewPanel.tsx
│       │   │   ├── MessageList.tsx
│       │   │   └── InterviewInput.tsx
│       │   ├── notebook/
│       │   │   ├── Notebook.tsx
│       │   │   ├── EvidenceTab.tsx
│       │   │   ├── ClaimsTab.tsx
│       │   │   ├── PeopleTab.tsx
│       │   │   └── TimelineTab.tsx
│       │   ├── correction/
│       │   │   └── PublishCorrectionDialog.tsx
│       │   └── endgame/
│       │       ├── SubmitTruthDialog.tsx
│       │       ├── GameWonDialog.tsx
│       │       └── GameLostDialog.tsx
│       ├── stores/
│       │   ├── gameStore.ts
│       │   └── uiStore.ts
│       ├── lib/
│       │   ├── api.ts
│       │   ├── ws.ts
│       │   ├── schemas.ts
│       │   └── constants.ts
│       └── types/
│           └── wire.ts
│
├── scripts/
│   ├── dev.ps1
│   ├── dev.sh
│   ├── test.ps1
│   ├── test.sh
│   └── scenario-check.sh
│
└── .github/
    └── workflows/
        ├── ci.yml
        └── race.yml
```

# 24. FILE-BY-FILE RESPONSIBILITY CATALOG
A coding agent must be able to explain why every file exists. The following is the minimum intended responsibility map.
| File / path | Exact responsibility |
|---|---|
| `backend/cmd/api/main.go` | Composition root. Load config, logger, DB pool, repositories, services, simulation manager, WebSocket hub, HTTP router; handle graceful shutdown. |
| `backend/internal/config/config.go` | Typed Config structs + Load() + Validate(). Only place mapping environment variables to application settings. |
| `backend/internal/config/defaults.go` | All non-secret V1 defaults: ports, tick duration, belief thresholds, queue sizes, timeouts. |
| `backend/internal/domain/session.go` | GameSession aggregate: status, scenario id, game time, player credibility, outcome; invariant methods. |
| `backend/internal/domain/agent.go` | Agent persistent identity/state; trait validation; no HTTP/DB imports. |
| `backend/internal/domain/claim.go` | Claim and ClaimClass; server-only truth classification. |
| `backend/internal/domain/belief.go` | AgentBelief value model, adoption helpers, clamp logic. |
| `backend/internal/domain/memory.go` | AgentMemory types and provenance. |
| `backend/internal/domain/evidence.go` | EvidenceItem, EvidenceLink, discovery requirements. |
| `backend/internal/domain/relationship.go` | Directed trust relation. |
| `backend/internal/domain/location.go` | Scenario location + waypoint position metadata. |
| `backend/internal/domain/world_event.go` | Canonical event envelope persisted and broadcast. |
| `backend/internal/domain/player_investigation.go` | Discovered evidence, interview flags, hypothesis/submission state. |
| `backend/internal/domain/errors.go` | Sentinel/domain errors safe for application mapping. |
| `backend/internal/domain/ids.go` | Typed string IDs or aliases + validation helpers. |
| `backend/internal/game/belief/propagation.go` | Deterministic propagation formula. Pure function first; no DB. |
| `backend/internal/game/belief/narrative.go` | Aggregate several claims into PrimaryFalseNarrative adoption ratio. |
| `backend/internal/game/belief/correction.go` | Public correction impact calculation and credibility adjustment. |
| `backend/internal/game/evidence/strength.go` | Coverage/reliability/diversity scoring. |
| `backend/internal/game/outcome/evaluator.go` | Win/loss evaluator; loss check has priority after state-changing event. |
| `backend/internal/game/clock/clock.go` | Game clock abstraction; production/tunable test clock. |
| `backend/internal/game/scheduler/scheduler.go` | Schedules autonomous NPC interactions without LLM blocking. |
| `backend/internal/game/scheduler/policy.go` | Candidate selection rules, cooldowns, conversation frequencies, seeded RNG. |
| `backend/internal/game/movement/waypoints.go` | Simple route graph and movement interpolation state. |
| `backend/internal/game/dialogue/context_builder.go` | Builds privacy-safe, knowledge-bounded prompt context for one NPC. |
| `backend/internal/game/dialogue/guard.go` | Rejects LLM output referencing unknown claims/evidence or violating output schema. |
| `backend/internal/game/dialogue/template_provider.go` | LLM-off deterministic dialogue provider. |
| `backend/internal/game/dialogue/openai_compatible.go` | Optional OpenAI-compatible HTTP provider; no vendor lock-in. |
| `backend/internal/game/dialogue/queue.go` | Bounded async dialogue work queue and worker pool. |
| `backend/internal/game/simulation/engine.go` | One session simulation engine; tick orchestration only. |
| `backend/internal/game/simulation/tick.go` | Tick step ordering and timing measurements. |
| `backend/internal/game/simulation/manager.go` | Lifecycle for multiple sessions; start/stop/recover active sessions. |
| `backend/internal/game/scenario/loader.go` | Read scenario JSON bundle. |
| `backend/internal/game/scenario/validator.go` | Cross-file referential/invariant validation. |
| `backend/internal/application/ports/repositories.go` | Repository interfaces used by services. |
| `backend/internal/application/ports/realtime.go` | Realtime publisher interface. |
| `backend/internal/application/ports/dialogue.go` | Dialogue provider interface. |
| `backend/internal/application/services/start_session.go` | Create and seed a new session transactionally. |
| `backend/internal/application/services/get_state.go` | Build player-safe state DTO; never leak truth class/private memories. |
| `backend/internal/application/services/interview.go` | Validate interaction, append conversation, enqueue/generate response. |
| `backend/internal/application/services/inspect.go` | Inspect location/evidence discovery flow. |
| `backend/internal/application/services/publish_correction.go` | Validate evidence selection, calculate correction, emit events. |
| `backend/internal/application/services/submit_truth.go` | Score atomic answer and invoke outcome evaluator. |
| `backend/internal/application/services/notebook.go` | Build notebook projection from only discovered data. |
| `backend/internal/infrastructure/postgres/pool.go` | pgxpool creation, health checks and connection settings. |
| `backend/internal/infrastructure/postgres/tx.go` | Transaction helper with explicit isolation where needed. |
| `backend/internal/infrastructure/postgres/repository/*.go` | Concrete repository implementations wrapping generated sqlc queries. |
| `backend/internal/infrastructure/postgres/sqlc/*.go` | Generated code only; never hand edit. |
| `backend/internal/infrastructure/realtime/hub.go` | WebSocket connection registry per game session. |
| `backend/internal/infrastructure/realtime/client.go` | Per-client read/write pumps, bounded send queue, disconnect policy. |
| `backend/internal/infrastructure/realtime/messages.go` | Wire protocol envelopes and payload structs. |
| `backend/internal/infrastructure/realtime/publisher.go` | Implements application realtime port. |
| `backend/internal/http/router.go` | chi router, middleware order, route mounting. |
| `backend/internal/http/middleware/request_id.go` | Request correlation ID. |
| `backend/internal/http/middleware/recovery.go` | Panic recovery + structured logging. |
| `backend/internal/http/middleware/cors.go` | Development/production CORS policy. |
| `backend/internal/http/handlers/sessions.go` | Start/get session HTTP handlers. |
| `backend/internal/http/handlers/agents.go` | Agent list/details/interview HTTP handlers. |
| `backend/internal/http/handlers/investigation.go` | Inspect, evidence, notebook, correction, truth submission. |
| `backend/internal/http/handlers/health.go` | /health/live and /health/ready. |
| `backend/internal/http/handlers/ws.go` | Upgrade HTTP to WebSocket and join session hub. |
| `backend/internal/http/respond/json.go` | Consistent JSON success/error response helpers. |
| `backend/db/migrations/*.sql` | Versioned PostgreSQL DDL. |
| `backend/db/queries/*.sql` | Handwritten SQL annotated for sqlc. |
| `backend/sqlc.yaml` | sqlc generation config. |
| `backend/go.mod` | Go 1.27 module and dependencies. |
| `backend/go.sum` | Committed dependency lock hashes. |
| `frontend/src/app/page.tsx` | Landing/start investigation surface. |
| `frontend/src/app/game/[sessionId]/page.tsx` | Game route; fetch initial state and mount client game shell. |
| `frontend/src/components/game/GameShell.tsx` | Main orchestration layout for 3D canvas + HUD/panels. |
| `frontend/src/components/world/GameCanvas.tsx` | Client-only R3F Canvas. |
| `frontend/src/components/world/CityScene.tsx` | Scene composition: lighting, district, agents, rumor visual layer. |
| `frontend/src/components/world/AgentAvatar.tsx` | Clickable agent visual; interpolation only. |
| `frontend/src/components/world/LocationBuilding.tsx` | Clickable location anchor/building. |
| `frontend/src/components/world/CameraRig.tsx` | Orthographic camera, pan/zoom/limited rotation. |
| `frontend/src/components/world/RumorPulse.tsx` | Ephemeral visual cue for communication/correction events. |
| `frontend/src/components/hud/BeliefMeter.tsx` | Always-visible loss meter based only on server aggregate. |
| `frontend/src/components/hud/EvidenceMeter.tsx` | Investigation progress indicator. |
| `frontend/src/components/hud/EventFeed.tsx` | Player-safe event timeline. |
| `frontend/src/components/interview/InterviewPanel.tsx` | Free-text interview UX. |
| `frontend/src/components/notebook/Notebook.tsx` | Evidence/Claims/People/Timeline tabs. |
| `frontend/src/components/correction/PublishCorrectionDialog.tsx` | Select challenged claim + evidence + optional wording. |
| `frontend/src/components/endgame/SubmitTruthDialog.tsx` | Atomic ground-truth form. |
| `frontend/src/stores/gameStore.ts` | Server projections only; no business calculations. |
| `frontend/src/stores/uiStore.ts` | Local panel/camera selection state. |
| `frontend/src/lib/api.ts` | Typed fetch wrapper. |
| `frontend/src/lib/ws.ts` | WebSocket client with reconnect and sequence handling. |
| `frontend/src/lib/schemas.ts` | Zod schemas for critical wire payloads. |
| `frontend/src/types/wire.ts` | Transport types shared across frontend feature modules. |
| `frontend/public/models/*.glb` | Optimized low-poly assets. |
| `data/scenarios/riverside-factory/*.json` | Scenario source data loaded by backend. |
| `docs/GAME_BIBLE.md` | Player rules and design principles. |
| `docs/SCENARIO_RIVERSIDE.md` | Narrative/script data and ground-truth design. |
| `docs/ARCHITECTURE.md` | Runtime boundaries and diagrams. |
| `docs/PROTOCOL.md` | REST/WebSocket contracts. |
| `docs/ADR/*.md` | Architecture decision records. |
| `docs/IMPLEMENTATION_STATUS.md` | Coding-agent checkpoint ledger. |


# 25. DATABASE DESIGN

## 25.1 General rules

- PostgreSQL 18.x for V1 local/dev deployment.
- UUID primary keys for session/runtime entities.
- Scenario definition IDs remain stable strings (`agent_firefighter`, `claim_multiple_deaths`) and are copied into session rows where useful.
- `timestamptz` for wall-clock timestamps.
- integer `game_second` for simulation time.
- numeric probabilities stored as `double precision` with check constraints `[0,1]`.
- JSONB only for genuinely flexible payloads such as world-event payloads; do not hide normal relational data in JSONB.

## 25.2 Core tables

### `game_sessions`

```sql
create table game_sessions (
    id uuid primary key,
    scenario_id text not null,
    status text not null,
    simulation_seed bigint not null,
    game_second bigint not null default 0,
    player_credibility double precision not null
        check (player_credibility between 0 and 1),
    false_narrative_ratio double precision not null default 0
        check (false_narrative_ratio between 0 and 1),
    evidence_strength double precision not null default 0
        check (evidence_strength between 0 and 1),
    revision bigint not null default 0,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    ended_at timestamptz null
);
```

### `session_agents`

```sql
create table session_agents (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    scenario_agent_id text not null,
    name text not null,
    role text not null,
    current_location_id text not null,
    influence double precision not null check (influence between 0 and 1),
    credibility double precision not null check (credibility between 0 and 1),
    skepticism double precision not null check (skepticism between 0 and 1),
    sociality double precision not null check (sociality between 0 and 1),
    deception_tendency double precision not null check (deception_tendency between 0 and 1),
    receptiveness double precision not null check (receptiveness between 0 and 1),
    busy_until_game_second bigint not null default 0,
    revision bigint not null default 0,
    unique(session_id, scenario_agent_id)
);
```

### `session_claims`

Store a session copy so scenario files can evolve later without rewriting an existing saved session.

```sql
create table session_claims (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    scenario_claim_id text not null,
    display_text text not null,
    classification text not null,
    narrative_role text not null,
    narrative_weight double precision not null default 0,
    unique(session_id, scenario_claim_id)
);
```

`classification` never leaves server-safe projections.

### `agent_beliefs`

```sql
create table agent_beliefs (
    session_id uuid not null references game_sessions(id) on delete cascade,
    agent_id uuid not null references session_agents(id) on delete cascade,
    claim_id uuid not null references session_claims(id) on delete cascade,
    confidence double precision not null check (confidence between 0 and 1),
    first_source_agent_id uuid null references session_agents(id),
    last_source_agent_id uuid null references session_agents(id),
    root_source_agent_id uuid null references session_agents(id),
    evidence_resistance double precision not null default 0,
    updated_at_game_second bigint not null,
    revision bigint not null default 0,
    primary key(agent_id, claim_id)
);
```

### `agent_memories`

```sql
create table agent_memories (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    agent_id uuid not null references session_agents(id) on delete cascade,
    memory_type text not null,
    content text not null,
    reliability double precision not null check (reliability between 0 and 1),
    importance double precision not null check (importance between 0 and 1),
    source_agent_id uuid null references session_agents(id),
    root_source_agent_id uuid null references session_agents(id),
    related_claim_ids uuid[] not null default '{}',
    created_at_game_second bigint not null
);
```

### `relationships`

```sql
create table relationships (
    session_id uuid not null references game_sessions(id) on delete cascade,
    from_agent_id uuid not null references session_agents(id) on delete cascade,
    to_agent_id uuid not null references session_agents(id) on delete cascade,
    trust double precision not null check (trust between 0 and 1),
    revision bigint not null default 0,
    primary key(from_agent_id, to_agent_id)
);
```

### `evidence_items`

```sql
create table evidence_items (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    scenario_evidence_id text not null,
    name text not null,
    description text not null,
    location_id text not null,
    reliability double precision not null check (reliability between 0 and 1),
    source_category text not null,
    discovered boolean not null default false,
    discovered_at_game_second bigint null,
    unique(session_id, scenario_evidence_id)
);
```

### `evidence_claim_links`

```sql
create table evidence_claim_links (
    evidence_id uuid not null references evidence_items(id) on delete cascade,
    claim_id uuid not null references session_claims(id) on delete cascade,
    relation text not null check (relation in ('supports','contradicts')),
    weight double precision not null default 1 check (weight between 0 and 1),
    primary key(evidence_id, claim_id, relation)
);
```

### `conversations`

```sql
create table conversations (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    agent_id uuid not null references session_agents(id),
    started_at_game_second bigint not null,
    ended_at_game_second bigint null
);
```

### `conversation_messages`

```sql
create table conversation_messages (
    id uuid primary key,
    conversation_id uuid not null references conversations(id) on delete cascade,
    speaker_type text not null check (speaker_type in ('player','agent','system')),
    content text not null,
    intent text null,
    created_at_game_second bigint not null,
    model_provider text null,
    model_name text null,
    latency_ms integer null
);
```

### `player_discovered_evidence`

Could be derived from `evidence_items.discovered` in single-player V1. Keep a separate table only if multi-player or multi-investigator expansion requires per-player state. For V1, do not duplicate it.

### `publications`

```sql
create table publications (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    challenged_claim_id uuid not null references session_claims(id),
    message text not null,
    strength double precision not null check (strength between 0 and 1),
    player_credibility_before double precision not null,
    player_credibility_after double precision not null,
    published_at_game_second bigint not null
);
```

### `publication_evidence`

```sql
create table publication_evidence (
    publication_id uuid not null references publications(id) on delete cascade,
    evidence_id uuid not null references evidence_items(id),
    primary key(publication_id, evidence_id)
);
```

### `world_events`

```sql
create table world_events (
    id uuid primary key,
    session_id uuid not null references game_sessions(id) on delete cascade,
    sequence bigint not null,
    event_type text not null,
    game_second bigint not null,
    actor_agent_id uuid null references session_agents(id),
    target_agent_id uuid null references session_agents(id),
    payload jsonb not null default '{}'::jsonb,
    occurred_at timestamptz not null default now(),
    unique(session_id, sequence)
);
```

## 25.3 Essential indexes

```sql
create index idx_world_events_session_sequence
    on world_events(session_id, sequence);

create index idx_agent_memories_agent_time
    on agent_memories(agent_id, created_at_game_second desc);

create index idx_session_agents_session_location
    on session_agents(session_id, current_location_id);

create index idx_beliefs_session_claim_confidence
    on agent_beliefs(session_id, claim_id, confidence);

create index idx_evidence_session_discovered
    on evidence_items(session_id, discovered);
```

Do not add indexes without a query/use case.

---

# 26. TRANSACTION BOUNDARIES

## Start session

One DB transaction must:
1. create session,
2. copy scenario claims,
3. create 20 session agents,
4. create evidence rows,
5. create relationships,
6. seed observations/memories,
7. seed initial beliefs,
8. append `session.started`,
9. commit.

Only after commit does the simulation manager start ticking.

## Belief propagation

One propagation operation should transactionally:
1. read/check latest relevant belief revisions if needed,
2. insert memory for listener,
3. update listener belief,
4. append `rumor.shared`,
5. append `belief.changed`,
6. update aggregate ratio if affected,
7. evaluate outcome,
8. commit,
9. publish WebSocket events after commit.

Never broadcast an event for a DB transaction that later rolls back.

## Publish correction

One transaction:
1. validate all evidence belongs to session and is discovered,
2. compute strength,
3. create publication,
4. create publication-evidence links,
5. apply affected belief updates,
6. update player credibility,
7. recalc narrative ratio,
8. evaluate outcome,
9. append events,
10. commit.

---

# 27. HTTP API CONTRACT

Base:

```text
/api/v1
```

## `POST /api/v1/sessions`

Request:

```json
{
  "scenario_id": "riverside-factory",
  "seed": 424242
}
```

`seed` is optional in normal play and useful for replay/testing.

Response `201`:

```json
{
  "id": "uuid",
  "scenario_id": "riverside-factory",
  "status": "running",
  "game_second": 0,
  "primary_false_narrative_ratio": 0.05,
  "evidence_strength": 0.0,
  "player_credibility": 0.5,
  "revision": 1
}
```

## `GET /api/v1/sessions/{id}`

Returns player-safe snapshot.

Never includes server claim classification.

## `GET /api/v1/sessions/{id}/agents`

Player-safe list:

```json
[
  {
    "id": "uuid",
    "scenario_agent_id": "agent_security_guard",
    "name": "Minh Trần",
    "role": "Bảo vệ nhà máy",
    "location_id": "factory_gate",
    "public_status": "available",
    "interviewed": false
  }
]
```

Do not return exact skepticism/deception/private motive in normal gameplay.

## `GET /api/v1/sessions/{id}/agents/{agentId}`

Returns public profile + player-discovered notes.

## `POST /api/v1/sessions/{id}/agents/{agentId}/interviews`

Request:

```json
{
  "message": "Anh đã trực tiếp nhìn thấy điều gì lúc sự cố bắt đầu?"
}
```

Response can be:
- synchronous template answer when LLM off,
- accepted/pending with WebSocket result when async LLM is on.

Preferred V1 behavior: handler creates player message and returns `202 Accepted`; final agent response arrives by WebSocket.

## `POST /api/v1/sessions/{id}/locations/{locationId}/inspect`

Returns newly discovered items / narrative description.

## `GET /api/v1/sessions/{id}/notebook`

Returns:
- discovered evidence,
- player-known claims,
- interview status,
- player-safe timeline.

## `POST /api/v1/sessions/{id}/corrections`

Request:

```json
{
  "challenged_claim_id": "claim_multiple_deaths",
  "evidence_ids": [
    "ev_hospital_summary",
    "ev_dispatch_log"
  ],
  "message": "Bệnh viện xác nhận bốn ca nhẹ và không có tử vong."
}
```

## `POST /api/v1/sessions/{id}/truth-submissions`

Request:

```json
{
  "event_type": "fire",
  "cause_category": "electrical_fault",
  "major_explosion": false,
  "fatalities": 0,
  "chemical_release": false,
  "management_issue": "delayed_maintenance"
}
```

Response includes:
- correctness by allowed player-facing field,
- current evidence strength,
- terminal status if won,
- missing proof message when truth is right but evidence is insufficient.

Do not reveal correct answers for wrong submissions.

## Health

```text
GET /health/live
GET /health/ready
```

---

# 28. WEBSOCKET PROTOCOL

Endpoint:

```text
GET /ws?session_id=<uuid>
```

Client sends a resume message after connect:

```json
{
  "type": "client.resume",
  "last_sequence": 180
}
```

Server may replay a bounded event window or require snapshot refresh.

## Envelope

```json
{
  "sequence": 181,
  "type": "rumor.shared",
  "session_id": "uuid",
  "game_second": 901,
  "payload": {}
}
```

## Player-safe event types

```text
session.snapshot
session.status_changed
agent.moved
agent.public_status_changed
conversation.agent_message
rumor.visualized
belief_stats.updated
evidence.discovered
correction.published
player.credibility_updated
game.won
game.lost
server.resync_required
```

`rumor.visualized` must not necessarily reveal exact hidden claim content if the player has not learned it. It can show that “Minh talked to Thảo” or a category/icon.

## Backpressure

Each WebSocket client gets a bounded send queue.

Default:

```text
ws_send_queue = 256
```

If a slow client cannot keep up:
1. log structured warning,
2. close with an appropriate reason,
3. client reconnects and resyncs.

Never allow one slow browser to block a simulation session.

---

# 29. FRONTEND TECHNOLOGY

Baseline:
- Next.js 16.3.x Active LTS.
- React 19.
- TypeScript strict.
- `@react-three/fiber` 9.x to pair with React 19.
- Three.js.
- `@react-three/drei`.
- Zustand.
- Zod.
- Tailwind CSS.
- native WebSocket client.
- Lucide or similarly lightweight icon set.

Do not use SignalR because backend is Go.

---

# 30. FRONTEND STATE OWNERSHIP

## `gameStore`

Server-derived:
- session status,
- game second,
- false narrative ratio,
- evidence strength,
- player credibility,
- agent public projections,
- latest sequence,
- recent player-safe events.

## `uiStore`

Client-only:
- selected agent,
- selected location,
- notebook open,
- notebook tab,
- interview panel open,
- camera focus target,
- accessibility pause preference,
- local visual settings.

Never put propagation formulas in Zustand.

---

# 31. 3D WORLD IMPLEMENTATION

## 31.1 Agent representation

V1 may use:
- simple low-poly stylized human models,
- 4–6 base body variants,
- color/accessory variation by occupation.

Do not require unique bespoke character models for all 20 NPCs.

## 31.2 Animation states

Minimum:

```text
idle
walk
talk
inspect/working
```

## 31.3 Movement

Use waypoint graph.

Server owns semantic movement:
- from location A,
- to location B,
- departure game second,
- arrival game second.

Client interpolates visual position.

Do not send 60 position updates/sec.

Possible event:

```json
{
  "type": "agent.moved",
  "payload": {
    "agent_id": "uuid",
    "from": "cafe",
    "to": "market",
    "depart_game_second": 900,
    "arrive_game_second": 934
  }
}
```

Frontend derives smooth motion.

## 31.4 Rumor visual

When an NPC-to-NPC information event occurs and is player-visible:

- brief dotted arc/pulse,
- small speech icon,
- no giant laser beam,
- duration 0.8–1.5 seconds,
- repeated rapid events should be throttled visually.

The world should feel socially alive, not like a network-monitor dashboard.

---

# 32. UI SCRIPT

## Landing

Headline:

```text
CITY OF LIES
Find the truth before the lie becomes reality.
```

Primary CTA:

```text
Start Investigation
```

Secondary:
- “How it works”
- scenario card.

## Game HUD

Top:

```text
CITY OF LIES | 18:14 | FALSE BELIEF 34% | EVIDENCE 22% | CREDIBILITY 50%
```

Left/bottom:
- event feed.

Right:
- selected NPC/location panel.

Bottom action dock:
- Notebook
- Publish Correction
- Submit Truth

## Threshold feedback

`< 50%`
- normal.

`50–59%`
- subtle warning.

`60–69%`
- orange warning, social activity feels more intense.

`70–74.9%`
- critical red state, meter pulses but respect reduced-motion setting.

`>= 75%`
- terminal loss.

Do not reveal forecast like “you will lose in 2 minutes”.

---

# 33. NOTEBOOK UX

Tabs:

## Evidence

Each item:
- name,
- source type,
- reliability represented qualitatively initially (`weak / moderate / strong`) unless explicit numeric display is desired,
- what it supports/contradicts based on player interpretation UI.

## Claims

Player-known claims only.

Possible statuses:

```text
Unverified
Supported
Contradicted
Conflicting
```

These are notebook assessments, not server truth labels.

## People

- name,
- role,
- interviewed?,
- player note,
- known contradictions.

## Timeline

Only events the player has learned or public events.

Do not show hidden event history.

---

# 34. SCENARIO JSON CONTRACT

## `scenario.json`

```json
{
  "id": "riverside-factory",
  "title": "Riverside Factory Incident",
  "version": 1,
  "primary_false_narrative": {
    "adoption_threshold": 0.65,
    "defeat_ratio": 0.75,
    "claims": [
      {"claim_id": "claim_chemical_explosion", "weight": 0.35},
      {"claim_id": "claim_multiple_deaths", "weight": 0.40},
      {"claim_id": "claim_company_coverup", "weight": 0.25}
    ]
  },
  "required_truth_fields": [
    "event_type",
    "cause_category",
    "major_explosion",
    "fatalities"
  ]
}
```

## `truth-form.json`

```json
{
  "event_type": {
    "type": "enum",
    "answer": "fire",
    "options": ["fire", "chemical_explosion", "gas_explosion", "unknown"]
  },
  "cause_category": {
    "type": "enum",
    "answer": "electrical_fault",
    "options": ["electrical_fault", "chemical_reaction", "arson", "unknown"]
  },
  "major_explosion": {
    "type": "boolean",
    "answer": false
  },
  "fatalities": {
    "type": "integer",
    "answer": 0,
    "min": 0,
    "max": 20
  }
}
```

The backend loader may parse answers into server-only truth structures that are never serialized to player DTOs.

---

# 35. ENVIRONMENT VARIABLES

Root `.env.example`:

```dotenv
# Database
POSTGRES_DB=cityoflies
POSTGRES_USER=cityoflies
POSTGRES_PASSWORD=cityoflies_dev_only
POSTGRES_PORT=5432

# Backend
CITYOFLIES_HTTP_ADDR=:8080
CITYOFLIES_DATABASE_URL=postgres://cityoflies:cityoflies_dev_only@postgres:5432/cityoflies?sslmode=disable
CITYOFLIES_SCENARIO_DIR=/app/data/scenarios

# Game
CITYOFLIES_DEFEAT_FALSE_BELIEF_RATIO=0.75
CITYOFLIES_BELIEF_ADOPTION_THRESHOLD=0.65
CITYOFLIES_TRUTH_EVIDENCE_THRESHOLD=0.70
CITYOFLIES_STARTING_PLAYER_CREDIBILITY=0.50

# Simulation
CITYOFLIES_TICK_MS=1000
CITYOFLIES_BELIEF_UPDATE_RATE=0.38
CITYOFLIES_DIRECT_EVIDENCE_MULTIPLIER=0.35
CITYOFLIES_SAME_ROOT_SOURCE_MULTIPLIER=0.45
CITYOFLIES_DEFAULT_RELATIONSHIP_TRUST=0.50
CITYOFLIES_MAX_NEW_CONVERSATIONS_PER_TICK=2
CITYOFLIES_PAIR_COOLDOWN_SECONDS=90

# Realtime
CITYOFLIES_WS_SEND_QUEUE=256
CITYOFLIES_WS_WRITE_TIMEOUT_MS=5000
CITYOFLIES_WS_PING_INTERVAL_MS=20000

# LLM
CITYOFLIES_LLM_ENABLED=false
CITYOFLIES_LLM_PROVIDER=template
CITYOFLIES_LLM_BASE_URL=
CITYOFLIES_LLM_API_KEY=
CITYOFLIES_LLM_MODEL=
CITYOFLIES_LLM_WORKERS=3
CITYOFLIES_LLM_QUEUE_CAPACITY=50
CITYOFLIES_LLM_TIMEOUT_MS=15000
CITYOFLIES_LLM_RETRY_COUNT=1

# Frontend
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
```

Production overrides secrets via hosting environment, not committed `.env`.

---

# 36. GO SETTINGS AND TOOLING

`go.mod`:

```text
go 1.27
```

Use a current Go 1.27 patch release in CI/container.

Recommended direct dependency families:
- `github.com/go-chi/chi/v5`
- `github.com/coder/websocket`
- `github.com/jackc/pgx/v5`
- UUID library if needed
- generated sqlc package code
- testcontainers only under tests if selected

Use:
- `go fmt`
- `go vet`
- `go test`
- `go test -race`
- `staticcheck` or golangci-lint configuration, but do not enable hundreds of noisy linters just to satisfy a badge.

---

# 37. TYPESCRIPT SETTINGS

Important `tsconfig` intent:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noImplicitOverride": true,
    "forceConsistentCasingInFileNames": true
  }
}
```

Do not set `skipLibCheck: false` purely for ideology if third-party types make builds impractical; document whatever Next-generated config requires.

No unbounded `any`.

At HTTP/WebSocket boundaries use Zod for critical message validation.

---

# 38. DOCKER COMPOSE

V1 services:

```text
postgres
backend
frontend
```

No Redis.

No message broker.

No separate AI microservice.

### PostgreSQL

Pin a stable PostgreSQL 18 tag appropriate to the environment rather than `latest`.

### Backend

Multi-stage:
1. Go build image,
2. runtime image,
3. copy binary + scenario data,
4. non-root user if practical.

### Frontend

Use Next standalone output if configured.

### Health dependency

- postgres has `pg_isready`,
- backend `/health/ready`,
- frontend depends on backend readiness only for compose orchestration, while browser handles transient startup failures gracefully.

---

# 39. LOGGING

Use JSON structured logs in production.

Fields:

```text
level
msg
request_id
session_id
agent_id
event_type
game_second
sequence
duration_ms
error
```

Never log:
- API keys,
- full LLM prompt with hidden/private state in production,
- sensitive system prompt,
- raw DB DSN password.

Development can log redacted context summaries.

---

# 40. METRICS

V1 development metrics worth exposing:

```text
active_sessions
simulation_tick_duration_ms
simulation_tick_overrun_total
npc_conversation_total
belief_change_total
websocket_clients
websocket_dropped_clients_total
llm_queue_depth
llm_request_duration_ms
llm_request_failure_total
```

Metrics are observability, not player UI.

---

# 41. TEST PLAN

## 41.1 Pure belief tests

Must include:

1. confidence never exceeds 1,
2. confidence never below 0,
3. high-trust source has greater delta than low-trust equivalent,
4. skeptical listener receives smaller delta,
5. contradictory direct observation reduces delta,
6. same-root repeated rumor has smaller effect,
7. deterministic formula returns exact expected value within tolerance.

## 41.2 Narrative tests

- weights aggregate correctly,
- agent at `0.6499` not adopted for 0.65 threshold,
- agent at `0.65` adopted,
- 14/20 agents = 70%, no loss,
- 15/20 = 75%, loss,
- ratio denominator excludes ineligible/disabled NPC only when explicitly configured.

## 41.3 Evidence tests

- duplicate evidence does not increase source diversity twice,
- high reliability increases reliability score,
- irrelevant evidence excluded,
- missing required truth field prevents full coverage,
- strength clamps `[0,1]`.

## 41.4 Outcome tests

- correct truth + 0.69 evidence => no win,
- correct truth + 0.70 evidence + false ratio 0.74 => win,
- correct truth + evidence 1.0 but false ratio already 0.75 => loss remains terminal,
- wrong truth never wins.

## 41.5 Scenario validation

Fail when:
- duplicate IDs,
- missing referenced agent,
- missing claim,
- evidence links to unknown claim,
- narrative weights do not sum near 1,
- truth form missing required field,
- agent start location missing,
- observation owner missing,
- relationship self-loop if forbidden,
- invalid probability.

## 41.6 Integration tests

With real PostgreSQL/testcontainer:
- session initialization transaction,
- session persistence/reload,
- interview write,
- correction transaction,
- world-event sequence uniqueness,
- reconnect snapshot.

## 41.7 Simulation soak test

Run accelerated, LLM off:

```text
100 sessions × several seeds
```

Track:
- percentage naturally reaching loss with no player,
- median time to 75%,
- any deadlock/panic,
- any invariant violation.

Initial balancing target:

```text
No-player loss rate: > 90%
Median no-player defeat: 12–18 real-equivalent minutes
```

The player must have to act.

## 41.8 Race test

Run backend core with:

```bash
go test -race ./...
```

---

# 42. SCENARIO BALANCING TARGETS

These are not guarantees; they are tuning goals.

No player intervention:

```text
18:14     15–25%
18:20     30–45%
18:24     40–55%
18:30     55–65%
18:35+    critical path toward 75%
```

Strong player route:
- interview Lan or Minh,
- obtain visual/direct evidence,
- reach hospital,
- obtain fire report,
- publish correction with independent sources,
- leverage Ngọc or Liên,
- submit truth.

Should allow win around:
```text
12–16 minutes
```

Weak route:
- spend time with low-quality rumor sources,
- publish too early,
- fail to corroborate,
- should often lose.

There must be multiple viable routes; do not make one hidden “correct click order”.

---

# 43. GOLDEN PATH — EXAMPLE SUCCESSFUL PLAYTHROUGH

This is a testable reference run, not the only solution.

### Minute 0
Player enters.

False belief ~5–10%.

### Minute 1
Player interviews Minh:
> “Anh trực tiếp thấy gì?”

Receives:
- smoke from storage,
- small sharp bang,
- clue about CCTV.

### Minute 2
Player inspects security room.
Discovers CCTV.

### Minute 3
Player interviews Lan.
Learns sparks near DB-4.

### Minute 4
False belief ~20–30%.

Player moves to hospital.

### Minute 5
Interviews Mai.
Unlocks hospital summary: 0 deaths.

### Minute 6
Asks follow-up about ambulance volume.
Unlocks dispatch log.

### Minute 7
False belief ~35–45%.

Player moves to fire station.

### Minute 8
Interviews Quang using collected evidence.
Unlocks fire report.

### Minute 9
Evidence strength now crosses a useful threshold.

Player publishes correction against `multiple deaths`
using:
- hospital summary,
- dispatch log.

### Minute 10
False-belief growth slows.
Credibility increases.

### Minute 11
Player publishes/feeds fire report through Ngọc or Liên.

### Minute 12
Chemical-explosion belief begins decreasing in skeptical/high-credibility cluster.

### Minute 13
Player submits:
- fire,
- electrical fault,
- no major explosion,
- zero fatalities.

Evidence >= 70%, false narrative < 75%.

Game won.

---

# 44. FAILURE PLAYTHROUGH — EXAMPLE

### Minute 0–4
Player interviews Hùng, Thảo, Vy repeatedly.

Receives mostly rumor.

### Minute 5
Publishes:
> “Không có gì nghiêm trọng.”

But has no hospital or fire evidence.

Correction strength low.
Credibility drops.

### Minute 7–10
Vy/Khoa cluster continues amplification.

### Minute 11
Player spends time chasing “toxic cloud” secondary rumor.

### Minute 14
Primary false narrative reaches ~68%.

### Minute 15
Player finally gets hospital evidence but cannot distribute a sufficiently strong correction in time.

### Minute 16
15/20 agents cross adoption threshold.

False narrative ratio = 75%.

Game lost.

End report explains:
- too much time spent on low-value sources,
- first correction lacked independent evidence,
- player credibility was damaged,
- high-influence nodes were not countered.

---

# 45. AFTER-ACTION REPORT

Win or loss screen should show:

```text
Outcome
Time used
Final false-belief ratio
Final evidence strength
Player credibility
Ground Truth accuracy
Evidence coverage
Corrections published
Most influential rumor source
Most useful evidence
Critical decision timeline
```

Do not claim pedagogical/scientific accuracy beyond the simulation.

Include a button:

```text
Replay same seed
```

and:

```text
New seed
```

V1 can implement replay as restart with same seed, not full event playback UI.

---

# 46. IMPLEMENTATION PHASES

## PHASE 0 — Repository + documents

Create:
- root structure,
- this master spec split into docs,
- ADRs,
- implementation status.

Exit:
- no code architecture ambiguity.

## PHASE 1 — Go skeleton

Create:
- Go module,
- config,
- logger,
- health server,
- graceful shutdown,
- tests.

Exit:
```bash
go test ./...
go vet ./...
```

pass.

## PHASE 2 — Pure domain/game math

Implement:
- claims,
- beliefs,
- narrative,
- evidence score,
- outcome evaluator,
- seeded RNG abstraction.

No DB.

Exit:
all pure unit tests pass.

## PHASE 3 — Scenario loader

Load Riverside JSON.

Exit:
`scenario-check` command validates all files.

## PHASE 4 — PostgreSQL

Migrations + pgx + sqlc.

Exit:
integration test creates and loads a session.

## PHASE 5 — Simulation engine

Implement ticks, scheduler, propagation, event append.

Exit:
headless `cmd/sim` can run a full no-player session to loss deterministically.

This is a critical milestone.

## PHASE 6 — HTTP + WebSocket

Create API contracts and realtime hub.

Exit:
a tiny dev page/CLI can watch false-belief ratio evolve.

## PHASE 7 — Frontend application shell

Landing, session creation, stores, WebSocket, HUD.

No 3D required yet.

Exit:
browser can play minimal text-mode investigation.

## PHASE 8 — Full investigation mechanics

Interview template provider, inspect, evidence, notebook, correction, truth submission.

Exit:
complete game can be won/lost without LLM and without 3D.

This is the **gameplay-complete milestone**.

## PHASE 9 — 3D world

Add district, agents, movement, click interaction, rumor visuals.

Exit:
3D presentation does not change game rules.

## PHASE 10 — LLM integration

Add bounded provider.

Exit:
LLM-on and LLM-off both pass core acceptance tests.

## PHASE 11 — balancing/polish

Seeds, UX, performance, errors, after-action report.

## PHASE 12 — CV/release quality

README, screenshots, architecture diagram, Docker, CI, demo deployment.

---

# 47. DEFINITION OF DONE FOR V1

V1 is not done because the city renders.

It is done only if:

1. Fresh clone can be launched from documented steps.
2. Riverside scenario loads from data files.
3. Twenty NPCs exist.
4. NPC conversations happen without player input.
5. Beliefs change deterministically.
6. Rumor reaches loss state if player ignores it.
7. Player can interview.
8. Player can inspect.
9. Player can discover evidence.
10. Notebook works.
11. Player can publish corrections.
12. Corrections measurably alter beliefs.
13. Player can submit truth.
14. Correct truth without evidence does not win.
15. 75% false narrative causes immediate terminal loss.
16. Game survives browser refresh.
17. WebSocket reconnect resyncs.
18. LLM disabled mode remains fully playable.
19. Go tests pass.
20. Race tests pass for core backend.
21. Frontend typecheck/build passes.
22. Scenario validation passes.
23. No hidden Ground Truth reaches browser network payloads.
24. No API keys are committed.
25. README accurately describes what was actually implemented.

---

# 48. CODING-AGENT OPERATING INSTRUCTIONS

The following section can be pasted into a coding agent as the operating contract.

## Role

You are implementing City of Lies as a production-quality portfolio project.

You are not allowed to silently simplify core simulation rules.

## Before editing

1. Read `docs/GAME_BIBLE.md`.
2. Read `docs/ARCHITECTURE.md`.
3. Read `docs/IMPLEMENTATION_STATUS.md`.
4. Read relevant ADRs.
5. Inspect existing code before creating duplicate abstractions.

## During implementation

- Work one phase at a time.
- Keep repository compiling at every checkpoint.
- Write tests with domain code.
- Do not begin LLM integration before deterministic game is winnable/losable.
- Do not begin visual polish before text-mode gameplay works.
- Do not replace Go backend.
- Do not introduce Redis/message broker/microservices in V1.
- Do not move game truth into prompts.
- Do not hand-edit sqlc generated files.
- Do not hard-code scenario IDs inside reusable game rules.
- Do not expose developer fields in player DTOs.

## At end of every task

Run relevant commands.

Backend:

```bash
gofmt -w .
go vet ./...
go test ./...
```

When concurrency changed:

```bash
go test -race ./...
```

Frontend:

```bash
npm run lint
npm run typecheck
npm test
npm run build
```

Update:

```text
docs/IMPLEMENTATION_STATUS.md
```

with:

```text
Phase
Completed
Files changed
Tests run
Known issue
Next exact task
```

## If architecture must change

Create ADR first:

```text
Context
Problem
Options
Decision
Consequences
Migration impact
```

Then implement.

---

# 49. IMPLEMENTATION STATUS TEMPLATE

```md
# Implementation Status

## Current phase
PHASE X — ...

## Completed
- ...

## Files created/changed
- `...`

## Tests executed
- `go test ./...` — PASS
- ...

## Manual acceptance
- ...

## Known issues
- ...

## Architecture decisions pending
- None / ...

## Next exact task
1. ...
2. ...
```

---

# 50. README POSITIONING FOR CV

Do not describe this merely as:

> “A 3D AI game.”

Preferred:

> **City of Lies is a browser-based real-time multi-agent simulation in which autonomous characters form and propagate beliefs from partial observations, social trust and memory. Players investigate the same event and must establish the evidence-backed Ground Truth before a false narrative reaches social dominance.**

Engineering bullets become truthful only after implemented:

- Built a deterministic multi-agent belief propagation engine in Go.
- Designed authoritative server-side truth, memory, claim and evidence boundaries.
- Implemented asynchronous LLM dialogue with knowledge guards and deterministic fallbacks.
- Built WebSocket real-time synchronization with sequence-based resync.
- Used PostgreSQL + pgx + sqlc for persistent sessions and event history.
- Built a 3D browser visualization with Next.js and React Three Fiber.
- Added seeded simulation tests and race testing for concurrent session execution.

---

# 51. NON-GOALS / ANTI-SCOPE

Reject these until after V1:

- “Let’s add accounts first.”
- “Let’s add multiplayer.”
- “Let’s make every NPC use an LLM planner.”
- “Let’s add vector DB because AI.”
- “Let’s add Kafka for events.”
- “Let’s add K8s.”
- “Let’s generate infinite scenarios.”
- “Let’s add combat.”
- “Let’s build a realistic Hanoi map.”
- “Let’s add mobile joystick.”
- “Let’s give every NPC voice.”
- “Let’s add blockchain reputation.”
- “Let’s make the browser calculate beliefs.”

Every one of these can wait.

---

# 52. FUTURE V2 DIRECTIONS

Only after V1:

- scenario authoring tool,
- multiple misinformation narratives,
- competing true narratives/uncertain ground truth,
- procedural NPC networks,
- multiplayer journalist teams,
- adversarial “misinformation actor” role,
- more sophisticated memory decay,
- graph visualization for replay,
- social-platform channels,
- research/education modes,
- scenario benchmarking.

None should block V1.

---

# 53. MASTER ACCEPTANCE SCRIPT

A human tester should be able to execute this exact script:

1. Run local stack.
2. Open landing page.
3. Start Riverside Factory Incident.
4. See 20 NPC agents.
5. Wait 2 minutes without acting.
6. Observe social interactions and increasing false-belief meter.
7. Select Minh.
8. Interview Minh.
9. Discover or unlock CCTV lead.
10. Inspect CCTV.
11. Travel to hospital.
12. Interview Mai.
13. Obtain hospital summary.
14. Open Notebook.
15. Confirm evidence and learned claims are listed.
16. Publish a correction against “multiple workers died”.
17. Observe correction event and changed belief trend.
18. Obtain fire report.
19. Open Submit Truth.
20. Submit correct atomic truth.
21. If Evidence Strength >= 70% and false ratio < 75%, see Win.
22. Restart with same seed.
23. Take no action.
24. Confirm scenario can progress to 75% and Lose.
25. Restart with LLM disabled.
26. Confirm same complete gameplay remains possible.
27. Refresh browser mid-session.
28. Confirm session state restores.
29. Disconnect/reconnect WebSocket.
30. Confirm snapshot/resync works.

If any of steps 1–30 fails, V1 is not complete.

---

# 54. FINAL DESIGN PRINCIPLE

The project succeeds if the player can say:

> “I knew there was a lie spreading, but I did not know who to trust. Every person knew only a piece, and while I was checking one lead, everyone else kept talking.”

The engineering succeeds if the developer can say:

> “The agents are not just prompts. They are persistent stateful actors inside a deterministic simulation whose beliefs, memories, relationships, evidence exposure and communication are governed by explicit rules. The LLM only inhabits the character.”

That distinction is the identity of **City of Lies**.

# APPENDIX A — INITIAL OBSERVATION SEED
The following section defines suggested initial knowledge. Exact numeric values can be tuned, but the semantic ownership must remain.

## `agent_security_guard` — Minh Trần

- Direct observation summary: Thấy khói bốc từ kho điện; nghe tiếng nổ nhỏ kiểu aptomat/thiết bị chập, không phải nổ hóa chất; thấy 2 công nhân ho sặc khi chạy ra.
- Blind spots: Không vào bên trong kho; không biết báo cáo bệnh viện; không biết nguyên nhân kỹ thuật cuối cùng.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_factory_worker_a` — Lan Phạm

- Direct observation summary: Ở gần kho khi cháy; thấy tia lửa ở tủ điện; bị ho do khói; được sơ cứu, không thấy ai tử vong.
- Blind spots: Không biết số người bị thương toàn bộ; không biết lịch bảo trì.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_factory_worker_b` — Hùng Nguyễn

- Direct observation summary: Không chứng kiến trực tiếp; nghe Lan nói 'có tiếng nổ ở tủ điện' và biến thành 'có vụ nổ'.
- Blind spots: Hầu như toàn bộ sự thật kỹ thuật.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.
- This character is an early amplification candidate; initialize relevant false claims slightly above unknown baseline only when justified by the timeline.

## `agent_firefighter` — Quang Lê

- Direct observation summary: Đội của anh dập cháy tại kho; không ghi nhận dấu vết nổ hóa chất; kết luận sơ bộ điểm cháy gần tủ phân phối điện.
- Blind spots: Không có quyền kết luận pháp y cuối cùng; không biết động cơ của người lan tin.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_paramedic` — Mai Vũ

- Direct observation summary: Tiếp nhận 4 ca nhẹ: 2 hít khói, 1 trầy xước khi sơ tán, 1 hoảng loạn; không có tử vong.
- Blind spots: Không biết chính xác nguồn cháy.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_shop_owner` — Thảo Đỗ

- Direct observation summary: Nhìn thấy xe cứu hỏa và xe cấp cứu chạy qua; không thấy hiện trường.
- Blind spots: Không biết thương vong hay nguyên nhân.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_resident` — Bác Dũng

- Direct observation summary: Thấy cột khói xám, ngửi thấy mùi khét điện/nhựa; không ngửi thấy mùi hóa chất lạ.
- Blind spots: Không có thông tin bên trong nhà máy.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_blogger` — Khoa Bùi

- Direct observation summary: Nhận ảnh khói từ Nhiếp ảnh gia và tin nhắn 'nghe nói có nổ' từ Hùng.
- Blind spots: Không có tài liệu y tế hay cứu hỏa.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.
- This character is an early amplification candidate; initialize relevant false claims slightly above unknown baseline only when justified by the timeline.

## `agent_influencer` — Vy Hoàng

- Direct observation summary: Không chứng kiến; chỉ xem bài của Khoa và bình luận của cư dân.
- Blind spots: Mọi dữ kiện trực tiếp.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.
- This character is an early amplification candidate; initialize relevant false claims slightly above unknown baseline only when justified by the timeline.

## `agent_taxi_driver` — Tuấn Đặng

- Direct observation summary: Chở một công nhân rời khu vực; người này nói 'chỉ bị khói, chưa nghe ai chết'.
- Blind spots: Không biết người công nhân đó có đủ thông tin hay không.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_student` — An Nguyễn

- Direct observation summary: Không thấy hiện trường; có thể truy ra timestamp và phiên bản bài đăng khác nhau.
- Blind spots: Không biết dữ kiện hiện trường nếu chưa hỏi.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_journalist` — Ngọc Trịnh

- Direct observation summary: Biết cách liên hệ cứu hỏa/bệnh viện; ban đầu chưa có tài liệu.
- Blind spots: Chưa trực tiếp chứng kiến.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_factory_manager` — Sơn Vương

- Direct observation summary: Biết không có báo cáo tử vong; biết ticket bảo trì tủ điện đã quá hạn 12 ngày.
- Blind spots: Không tận mắt thấy thời điểm chập điện.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_technician` — Đức Phan

- Direct observation summary: Biết tủ DB-4 có lỗi quá nhiệt; đã tạo maintenance ticket; ảnh sau cháy cho thấy điểm hồ quang điện.
- Blind spots: Không biết ai bắt đầu tin đồn.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_nurse` — Hà Lương

- Direct observation summary: Biết không có tử vong và không có ca bỏng hóa chất.
- Blind spots: Không biết nguồn cháy.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_vendor` — Cô Hương

- Direct observation summary: Thấy xe cấp cứu đưa vài người đi nhưng tất cả đều tự đi hoặc ngồi được.
- Blind spots: Không biết tình trạng sau đó.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_retired_teacher` — Thầy Bình

- Direct observation summary: Không có quan sát trực tiếp đáng kể; có năng lực đánh giá mâu thuẫn giữa các lời kể.
- Blind spots: Thiếu dữ kiện gốc.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_delivery_driver` — Phúc Trần

- Direct observation summary: Đi ngang nhà máy trước cháy 10 phút và thấy xe kỹ thuật đỗ gần kho; không biết lý do.
- Blind spots: Không thấy thời điểm cháy.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_community_leader` — Cô Liên

- Direct observation summary: Có kênh liên hệ chính quyền/cứu hỏa nhưng phản hồi đến chậm.
- Blind spots: Không biết nội bộ nhà máy.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.

## `agent_photographer` — Tùng Lâm

- Direct observation summary: Chụp chuỗi ảnh timestamp: khói bắt đầu ở khu kho, không có fireball; ảnh cho thấy cửa kho còn nguyên kết cấu.
- Blind spots: Không biết tình trạng người bị thương.
- Initial false-narrative beliefs should remain low unless this character already received a rumor before session start.


# APPENDIX B — SUGGESTED FIRST 10 RUMOR TRANSMISSIONS

These are anchor interactions for seed/balancing tests. The scheduler can vary target timing after the first few anchors.

| Order | Speaker | Listener | Claim | Reason |
|---:|---|---|---|---|
| 1 | Lan | Hùng | ambiguous “there was a bang” memory | direct witness → social coworker |
| 2 | Hùng | Thảo | `claim_chemical_explosion` precursor | exaggeration |
| 3 | Thảo | Hương | explosion rumor | neighborhood hub |
| 4 | Hùng | Khoa | explosion + people trapped | blogger input |
| 5 | Tùng | Khoa | dramatic image | visual amplifier |
| 6 | Khoa | Vy | chemical explosion speculation | high-influence transfer |
| 7 | Vy | Tuấn | casualty speculation | public amplification |
| 8 | Thảo | Phúc | casualty speculation | market network |
| 9 | Vy | An | “people died” rumor | social channel |
| 10 | An | Khoa | challenge/question about source | potential corrective friction |

The system should be capable of the rumor evolving even if exact order 7–10 changes.

# APPENDIX C — DEVELOPER DIAGNOSTIC VIEW

Development-only diagnostics may show:

```text
Agent             ChemicalExplosion  MultipleDeaths  Coverup  Narrative
Hùng                   0.71              0.46         0.28      0.51
Khoa                   0.79              0.62         0.43      0.63
Vy                     0.83              0.71         0.55      0.71  ADOPTED
...
```

This screen is forbidden in normal gameplay.

Useful debug columns:

```text
agent
location
busy_until
current salient claim
confidence
root source
last source
memory count
next share cooldown
narrative score
adopted bool
```

# APPENDIX D — LOCAL COMMANDS TARGET

Root `Makefile` should eventually provide:

```bash
make dev
make backend
make frontend
make test
make test-race
make lint
make db-up
make db-down
make migrate-up
make sqlc
make scenario-check
make sim
```

Windows PowerShell scripts should mirror critical developer operations where convenient.

# APPENDIX E — SOURCE VERSION NOTES (2026-09-10)

At document generation time:

- Go 1.27.1 is a current stable Go 1.27 patch release.
- Next.js 16.3.3 is the current Active LTS security target advertised by the Next.js project.
- React Three Fiber 9.x is the branch paired with React 19.
- PostgreSQL 18.6 is a current stable PostgreSQL 18 release; PostgreSQL 19 is still beta at this date.
- pgx v5 remains the stable major line.

Pin exact versions in lock files/container tags when scaffolding. Re-check security advisories before deployment.


---

# APPENDIX F — OFFICIAL VERSION REFERENCES

These references were checked when this blueprint was generated on 2026-09-10.

- Go release history: https://go.dev/doc/devel/release
- Go 1.27 release notes: https://go.dev/doc/go1.27
- Next.js release/news page: https://nextjs.org/blog
- React Three Fiber installation/compatibility: https://r3f.docs.pmnd.rs/getting-started/installation
- PostgreSQL release page: https://www.postgresql.org/docs/release/
- pgx repository/documentation: https://github.com/jackc/pgx
- chi releases: https://github.com/go-chi/chi/releases
- coder/websocket repository: https://github.com/coder/websocket
- sqlc documentation: https://sqlc.dev/

When the repository is actually scaffolded, pin exact dependency versions in `go.mod`, `go.sum`, `package.json`, `package-lock.json`, and Docker image tags. Re-check security releases instead of blindly copying this document's patch versions.
