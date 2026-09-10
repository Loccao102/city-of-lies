# Scenario 01: Riverside Factory Incident

## Ground Truth
- **Time**: ~17:58
- **Event**: Electrical arcing in distribution panel DB-4 in the Riverside storage facility, followed by localized fire.
- **Casualties**: 4 minor medical treatments (smoke inhalation, abrasions, panic). 0 fatalities.
- **Explosion**: Sharp electrical pop/bang from circuit breaker trip, NOT a major chemical explosion.
- **Chemical Release**: No hazardous cloud.
- **Root Cause**: Overdue maintenance ticket on DB-4 (unresolved for 12 days).

## The False Narrative
> *"A chemical explosion happened at Riverside Factory. Several workers died, and management is hiding the real casualty count."*

### Component Claims & Weights:
- `claim_chemical_explosion` (Weight: 0.35)
- `claim_multiple_deaths` (Weight: 0.40)
- `claim_company_coverup` (Weight: 0.25)
- **Adoption Threshold**: $\ge 0.65$ per agent.
- **Defeat Threshold**: $\ge 75\%$ of eligible agents (15/20) adopting the false narrative.

## Cast of Characters (20 Agents)
1. `agent_security_guard` - Minh Trần (Bảo vệ)
2. `agent_factory_worker_a` - Lan Phạm (Công nhân kho)
3. `agent_factory_worker_b` - Hùng Nguyễn (Công nhân đóng gói)
4. `agent_firefighter` - Quang Lê (Lính cứu hỏa)
5. `agent_paramedic` - Mai Vũ (Nhân viên cấp cứu)
6. `agent_shop_owner` - Thảo Đỗ (Chủ quán tạp hóa)
7. `agent_resident` - Bác Dũng (Cư dân)
8. `agent_blogger` - Khoa Bùi (Blogger)
9. `agent_influencer` - Vy Hoàng (Influencer)
10. `agent_taxi_driver` - Tuấn Đặng (Tài xế taxi)
11. `agent_student` - An Nguyễn (Sinh viên)
12. `agent_journalist` - Ngọc Trịnh (Phóng viên)
13. `agent_factory_manager` - Sơn Vương (Quản lý)
14. `agent_technician` - Đức Phan (Kỹ thuật viên)
15. `agent_nurse` - Hà Lương (Y tá)
16. `agent_vendor` - Cô Hương (Người bán hàng rong)
17. `agent_retired_teacher` - Thầy Bình (Giáo viên nghỉ hưu)
18. `agent_delivery_driver` - Phúc Trần (Tài xế giao hàng)
19. `agent_community_leader` - Cô Liên (Trưởng nhóm cộng đồng)
20. `agent_photographer` - Tùng Lâm (Nhiếp ảnh gia)
