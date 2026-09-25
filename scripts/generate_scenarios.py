import json
import os

BASE_DIR = os.path.join(os.path.dirname(__file__), "..", "data", "scenarios")

def save_json(scenario_dir, filename, data):
    path = os.path.join(scenario_dir, filename)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"Saved {path}")

# ==============================================================================
# SCENARIO 2: metro-hospital-outbreak (Báo Động Bệnh Viện Metro)
# ==============================================================================
def create_hospital_scenario():
    sc_id = "metro-hospital-outbreak"
    s_dir = os.path.join(BASE_DIR, sc_id)

    metadata = {
        "id": sc_id,
        "title": "Báo Động Bệnh Viện Metro",
        "version": 1,
        "description": "Lúc 11:45 trưa, 18 y bác sĩ và bệnh nhân nhập viện cấp cứu do ngộ độc thực phẩm từ món cá ngừ hỏng tại căng tin bệnh viện. 0 ca tử vong. Tin đồn ác ý thổi phồng rằng mầm bệnh virus vũ khí sinh học đã rò rỉ từ phòng thí nghiệm ngầm.",
        "primary_false_narrative": {
            "adoption_threshold": 0.65,
            "defeat_ratio": 0.75,
            "claims": [
                {"claim_id": "claim_lab_virus_leak", "weight": 0.35},
                {"claim_id": "claim_patient_fatalities", "weight": 0.40},
                {"claim_id": "claim_cdc_quarantine_coverup", "weight": 0.25}
            ]
        },
        "required_truth_fields": ["event_type", "cause_category", "major_explosion", "fatalities"],
        "bonus_truth_fields": ["pathogen_type", "management_issue"]
    }

    locations = [
        {"id": "er_entrance", "name": "Cổng Khoa Cấp Cứu", "description": "Nơi xe cứu thương tấp nập đưa người nhập viện.", "position": {"x": -24, "y": 0, "z": -16}, "connected_locations": ["er_triage", "hospital_lobby", "pharmacy"]},
        {"id": "er_triage", "name": "Khu Phân Loại Bệnh", "description": "Bác sĩ kiểm tra sinh hiệu bệnh nhân ngộ độc.", "position": {"x": -32, "y": 0, "z": -24}, "connected_locations": ["er_entrance", "internal_ward", "icu"]},
        {"id": "hospital_lobby", "name": "Sảnh Chờ Bệnh Viện", "description": "Thân nhân bệnh nhân tụ tập lo lắng.", "position": {"x": -14, "y": 0, "z": -18}, "connected_locations": ["er_entrance", "canteen", "admin_office"]},
        {"id": "canteen", "name": "Căng Tin Bệnh Viện", "description": "Nơi cung cấp bữa trưa cá ngừ kho mặn cho ca trực sáng.", "position": {"x": -20, "y": 0, "z": -26}, "connected_locations": ["hospital_lobby", "nutrition_dept"]},
        {"id": "nutrition_dept", "name": "Kho Dinh Dưỡng", "description": "Khu bảo quản nguyên liệu thực phẩm và mẫu lưu 24h.", "position": {"x": -28, "y": 0, "z": -8}, "connected_locations": ["canteen"]},
        {"id": "admin_office", "name": "Văn Phòng Giám Đốc", "description": "Nơi ban giám đốc tổ chức họp khẩn ứng phó.", "position": {"x": 2, "y": 0, "z": 6}, "connected_locations": ["hospital_lobby", "lab_biosafety"]},
        {"id": "lab_biosafety", "name": "Phòng Xét Nghiệm Vi Sinh", "description": "Phòng xét nghiệm mẫu bệnh phẩm đạt chuẩn an toàn sinh học cấp 2.", "position": {"x": 26, "y": 0, "z": -20}, "connected_locations": ["admin_office", "icu"]},
        {"id": "icu", "name": "Khu Hồi Sức Tích Cực", "description": "Nơi điều trị các ca mất nước nặng cần truyền dịch.", "position": {"x": -6, "y": 0, "z": -6}, "connected_locations": ["er_triage", "lab_biosafety"]},
        {"id": "pharmacy", "name": "Nhà Thuốc Bệnh Viện", "description": "Quầy cấp phát điện giải và than hoạt tính khẩn cấp.", "position": {"x": 10, "y": 0, "z": -4}, "connected_locations": ["er_entrance", "street_cafe"]},
        {"id": "internal_ward", "name": "Khoa Nội Tiêu Hóa", "description": "Khoa tiếp nhận theo dõi bệnh nhân sau khi rửa ruột.", "position": {"x": -22, "y": 0, "z": 14}, "connected_locations": ["er_triage", "mortuary"]},
        {"id": "mortuary", "name": "Nhà Tang Lễ Bệnh Viện", "description": "Nhà xác bệnh viện, hiện hoàn toàn trống ca tử vong mới.", "position": {"x": -10, "y": 0, "z": 18}, "connected_locations": ["internal_ward"]},
        {"id": "street_cafe", "name": "Cà Phê Cổng Viện", "description": "Điểm tụ tập của cánh tài xế taxi và phóng viên hóng tin.", "position": {"x": 18, "y": 0, "z": 12}, "connected_locations": ["pharmacy", "hospital_lobby"]}
    ]

    claims = [
        {"id": "claim_poisoning_occurred", "text": "Có vụ ngộ độc thực phẩm tập thể tại bệnh viện", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_tuna_histamine_origin", "text": "Nguyên nhân do độc tố histamine từ lô cá ngừ bảo quản sai nhiệt độ", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_zero_fatalities", "text": "Không có bất kỳ ca tử vong nào, tất cả đều đang hồi phục", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_no_virus_leak", "text": "Không có mầm bệnh virus nguy hiểm nào bị rò rỉ", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_canteen_fridge_broken", "text": "Tủ đông căng tin đã hỏng cảm biến nhiệt 3 ngày trước", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_lab_virus_leak", "text": "Virus nhân tạo chết người rò rỉ từ phòng thí nghiệm ngầm", "classification": "false", "role": "primary_false", "weight": 0.35},
        {"id": "claim_patient_fatalities", "text": "Đã có 5 bệnh nhân co giật nôn ra máu và tử vong", "classification": "false", "role": "primary_false", "weight": 0.40},
        {"id": "claim_cdc_quarantine_coverup", "text": "Bệnh viện đang bị phong tỏa ngầm để bưng bít thảm họa dịch tễ", "classification": "false", "role": "primary_false", "weight": 0.25},
        {"id": "claim_black_bodybags_moved", "text": "Thấy xe tải chở túi đen rời khỏi nhà xác", "classification": "false", "role": "secondary_false", "weight": 0.0},
        {"id": "claim_rumor_panicking_staff", "text": "Nhân viên y tế hoang mang bỏ ca trực", "classification": "false", "role": "secondary_false", "weight": 0.0}
    ]

    agents = [
        {"id": "h_dr_nghiem", "name": "Bác sĩ Trưởng Nghiêm", "role": "Trưởng khoa Cấp cứu", "location_id": "er_triage", "influence": 0.65, "credibility": 0.88, "skepticism": 0.80, "sociality": 0.40, "deception_tendency": 0.05, "receptiveness": 0.45, "personality": "quyết đoán, chuyên nghiệp, nghiêm túc", "private_motive": "Khống chế tình hình cấp cứu nhanh nhất.", "direct_knowledge": "18 ca vào viện với triệu chứng dị ứng, đỏ bừng mặt, đau bụng tiêu chảy điển hình của ngộ độc histamine.", "blind_spots": "Không trực tiếp kiểm tra mẫu thức ăn ở căng tin."},
        {"id": "h_nurse_hoa", "name": "Y tá Hoa", "role": "Điều dưỡng cấp cứu", "location_id": "er_triage", "influence": 0.40, "credibility": 0.72, "skepticism": 0.50, "sociality": 0.75, "deception_tendency": 0.08, "receptiveness": 0.65, "personality": "chu đáo, dễ xúc động, nói nhiều", "private_motive": "Lo lắng cho sức khỏe các đồng nghiệp bị ngộ độc.", "direct_knowledge": "Tất cả bệnh nhân đều đã được truyền dịch và tiêm kháng histamine, mạch và SpO2 ổn định.", "blind_spots": "Nghe người ngoài đồn về phòng thí nghiệm sinh học."},
        {"id": "h_chef_tam", "name": "Bếp trưởng Tâm", "role": "Quản lý Căng tin", "location_id": "canteen", "influence": 0.35, "credibility": 0.60, "skepticism": 0.40, "sociality": 0.60, "deception_tendency": 0.25, "receptiveness": 0.55, "personality": "bồn chồn, lo sợ trách nhiệm", "private_motive": "Sợ bị tước giấy phép kinh doanh bếp ăn.", "direct_knowledge": "Lô cá ngừ mua sáng nay có dấu hiệu rã đông sớm nhưng vẫn cho chế biến.", "blind_spots": "Không hiểu cơ chế vi sinh vật."},
        {"id": "h_tech_minh", "name": "Kỹ thuật viên Minh", "role": "Phòng Vi sinh", "location_id": "lab_biosafety", "influence": 0.55, "credibility": 0.90, "skepticism": 0.85, "sociality": 0.30, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "chính xác, khoa học, ít biểu cảm", "private_motive": "Bảo vệ danh dự phòng thí nghiệm của viện.", "direct_knowledge": "Phòng lab chỉ lưu trữ chủng vi khuẩn thông thường để thử kháng sinh đồ, không có virus nguy hiểm cấp 4.", "blind_spots": "Không ra khu vực cấp cứu."},
        {"id": "h_guard_dung", "name": "Bảo vệ Dũng", "role": "Bảo vệ cổng viện", "location_id": "er_entrance", "influence": 0.35, "credibility": 0.58, "skepticism": 0.40, "sociality": 0.80, "deception_tendency": 0.15, "receptiveness": 0.70, "personality": "thích bàn tán, nhiều chuyện", "private_motive": "Tò mò, thích nghe ngóng các vụ việc giật gân.", "direct_knowledge": "Thấy nhiều bác sĩ mặc đồ bảo hộ đón xe cấp cứu nên tưởng dịch bệnh.", "blind_spots": "Không phân biệt được trang phục chống nhiễm khuẩn với đồ bảo hộ dịch."},
        {"id": "h_director_khanh", "name": "Giám đốc Khánh", "role": "Giám đốc Bệnh viện", "location_id": "admin_office", "influence": 0.80, "credibility": 0.82, "skepticism": 0.75, "sociality": 0.45, "deception_tendency": 0.10, "receptiveness": 0.40, "personality": "điềm tĩnh, cẩn trọng, ngoại giao", "private_motive": "Tránh để dư luận hoảng loạn gây hỗn loạn bệnh viện.", "direct_knowledge": "Đã có báo cáo nhanh từ khoa Dược và Cấp cứu xác nhận ngộ độc thực phẩm.", "blind_spots": "Chậm trễ họp báo giải thích cho người dân."},
        {"id": "h_reporter_nga", "name": "Phóng viên Nga", "role": "Nhà báo Y tế", "location_id": "street_cafe", "influence": 0.70, "credibility": 0.65, "skepticism": 0.60, "sociality": 0.85, "deception_tendency": 0.12, "receptiveness": 0.70, "personality": "sắc sảo, tìm kiếm tin tức độc quyền", "private_motive": "Muốn có bài phóng sự nóng nhất mạng xã hội.", "direct_knowledge": "Biết có nhiều xe cấp cứu đến liên tục và cửa viện đang hạn chế người ra vào.", "blind_spots": "Chưa tiếp cận được hồ sơ bệnh án nội bộ."},
        {"id": "h_relative_quang", "name": "Bác Quang", "role": "Thân nhân bệnh nhân", "location_id": "hospital_lobby", "influence": 0.38, "credibility": 0.50, "skepticism": 0.30, "sociality": 0.75, "deception_tendency": 0.05, "receptiveness": 0.85, "personality": "hoảng loạn, cả tin, lo sợ", "private_motive": "Lo cho con trai đang nằm truyền dịch trong phòng hồi sức.", "direct_knowledge": "Thấy con mình nổi ban đỏ khắp người và nôn thốc nôn tháo.", "blind_spots": "Tưởng triệu chứng nổi ban là do virus độc."},
        {"id": "h_pharmacist_yen", "name": "Dược sĩ Yến", "role": "Phụ trách Nhà thuốc", "location_id": "pharmacy", "influence": 0.45, "credibility": 0.82, "skepticism": 0.70, "sociality": 0.55, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "tỉ mỉ, nguyên tắc", "private_motive": "Đảm bảo cung cấp đủ thuốc giải độc đường tiêu hóa.", "direct_knowledge": "Khoa cấp cứu chỉ yêu cầu cấp Promethazine, than hoạt và Ringer Lactate, không có thuốc kháng virus đặc hiệu.", "blind_spots": "Không nắm số lượng bệnh nhân chính xác."},
        {"id": "h_inspector_hung", "name": "Thanh tra Hùng", "role": "Chi cục An toàn Vệ sinh Thực phẩm", "location_id": "nutrition_dept", "influence": 0.60, "credibility": 0.92, "skepticism": 0.88, "sociality": 0.35, "deception_tendency": 0.02, "receptiveness": 0.30, "personality": "khắt khe, công minh, điều tra kỹ", "private_motive": "Làm rõ nguồn gốc thực phẩm gây ngộ độc.", "direct_knowledge": "Đã niêm phong tủ lưu mẫu thức ăn trưa nay tại căng tin.", "blind_spots": "Chờ kết quả xét nghiệm hóa sinh định lượng."},
        {"id": "h_mortuary_van", "name": "Bác Vạn", "role": "Nhân viên Nhà tang lễ", "location_id": "mortuary", "influence": 0.30, "credibility": 0.75, "skepticism": 0.65, "sociality": 0.40, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "lầm lì, ít nói", "private_motive": "Chỉ làm tròn trách nhiệm.", "direct_knowledge": "Từ sáng tới giờ chưa có ca tử vong nào chuyển xuống nhà xác.", "blind_spots": "Không biết chuyện gì đang xảy ra ở sảnh chính."},
        {"id": "h_patient_tri", "name": "Bệnh nhân Trí", "role": "Bác sĩ nội trú bị ngộ độc", "location_id": "internal_ward", "influence": 0.50, "credibility": 0.80, "skepticism": 0.65, "sociality": 0.60, "deception_tendency": 0.05, "receptiveness": 0.55, "personality": "mệt mỏi nhưng còn tỉnh táo", "private_motive": "Muốn mau chóng khỏe lại để tiếp tục làm việc.", "direct_knowledge": "Ăn trưa cá ngừ lúc 11:30, đến 12:00 thì ngứa cổ họng, nóng bừng mặt và nôn ói.", "blind_spots": "Không theo dõi mạng xã hội đang đồn gì."},
        {"id": "h_blogger_khoa", "name": "Blogger Khoa", "role": "Tiktoker Reviewer", "location_id": "street_cafe", "influence": 0.75, "credibility": 0.35, "skepticism": 0.20, "sociality": 0.95, "deception_tendency": 0.30, "receptiveness": 0.80, "personality": "thích giật gân, câu like, phóng đại", "private_motive": "Tăng follow và lượt xem trên kênh cá nhân.", "direct_knowledge": "Chỉ quay video từ xa cảnh cổng viện đóng kín và xe hú còi liên tục.", "blind_spots": "Hoàn toàn bịa đặt về việc có người chết tím tái."},
        {"id": "h_cleaner_bich", "name": "Chị Bích", "role": "Tạp vụ bệnh viện", "location_id": "canteen", "influence": 0.25, "credibility": 0.55, "skepticism": 0.35, "sociality": 0.80, "deception_tendency": 0.10, "receptiveness": 0.75, "personality": "chất phác, hay chia sẻ chuyện hậu trường", "private_motive": "Muốn cảnh báo mọi người đừng ăn đồ căng tin.", "direct_knowledge": "Tủ lạnh căng tin bốc mùi chua từ hôm qua nhưng nhà bếp không chịu vứt cá.", "blind_spots": "Không biết chẩn đoán y khoa."},
        {"id": "h_driver_long", "name": "Tài xế Long", "role": "Lái xe cứu thương", "location_id": "er_entrance", "influence": 0.40, "credibility": 0.65, "skepticism": 0.50, "sociality": 0.75, "deception_tendency": 0.08, "receptiveness": 0.65, "personality": "nhiệt huyết, vội vã", "private_motive": "Chở bệnh nhân cấp cứu kịp thời.", "direct_knowledge": "Đón các bác sĩ bị đau bụng từ phân viện 2 sang, tất cả đều thở tốt.", "blind_spots": "Nghe đồn phong phanh có dịch lạ."},
        {"id": "h_student_linh", "name": "Sinh viên Linh", "role": "Thực tập sinh Y khoa", "location_id": "internal_ward", "influence": 0.35, "credibility": 0.70, "skepticism": 0.60, "sociality": 0.65, "deception_tendency": 0.05, "receptiveness": 0.60, "personality": "chăm chỉ, ham học hỏi", "private_motive": "Hỗ trợ các đàn anh xử lý ngộ độc.", "direct_knowledge": "Các bệnh nhân sau khi tiêm thuốc dị ứng đều đã giảm đỏ da và hết buồn nôn.", "blind_spots": "Không tiếp xúc thân nhân bên ngoài."},
        {"id": "h_officer_nam", "name": "Đại úy Nam", "role": "Công an khu vực", "location_id": "hospital_lobby", "influence": 0.60, "credibility": 0.85, "skepticism": 0.80, "sociality": 0.40, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "thận trọng, giữ gìn trật tự", "private_motive": "Ngăn chặn tụ tập đông người gây mất an ninh cổng viện.", "direct_knowledge": "Bệnh viện chỉ tạm dừng nhận bệnh thông thường để dồn sức cấp cứu hàng loạt.", "blind_spots": "Không có kiến thức dịch tễ chuyên sâu."},
        {"id": "h_nurse_phuong", "name": "Điều dưỡng Phương", "role": "Y tá Khoa Hồi sức", "location_id": "icu", "influence": 0.42, "credibility": 0.78, "skepticism": 0.65, "sociality": 0.50, "deception_tendency": 0.04, "receptiveness": 0.50, "personality": "tập trung, chịu áp lực tốt", "private_motive": "Theo dõi sát các chỉ số sinh tồn của bệnh nhân.", "direct_knowledge": "3 ca tụt huyết áp nhẹ đều đã ổn định lại sau khi bù 1000ml dịch muối sinh lý.", "blind_spots": "Không ra khỏi phòng hồi sức."},
        {"id": "h_vendor_mai", "name": "Cô Mai", "role": "Chủ quán nước cổng viện", "location_id": "street_cafe", "influence": 0.45, "credibility": 0.45, "skepticism": 0.25, "sociality": 0.90, "deception_tendency": 0.15, "receptiveness": 0.85, "personality": "buôn chuyện, bán nước kiêm trung tâm tin đồn", "private_motive": "Thu hút khách ghé uống nước nghe ngóng.", "direct_knowledge": "Thấy bảo vệ đóng cổng sắt ngăn người nhà xông vào.", "blind_spots": "Tự suy diễn là bệnh viện giấu dịch nguy hiểm."},
        {"id": "h_technician_son", "name": "Kỹ thuật viên Sơn", "role": "Bảo trì thiết bị lạnh", "location_id": "nutrition_dept", "influence": 0.32, "credibility": 0.75, "skepticism": 0.70, "sociality": 0.45, "deception_tendency": 0.15, "receptiveness": 0.45, "personality": "ít nói, cẩn thận", "private_motive": "Chứng minh sự cố tủ lạnh không phải do lỗi của mình.", "direct_knowledge": "Tủ đông căng tin bị đứt dây rơ-le nhiệt, nhiệt độ duy trì ở mức 15°C thay vì -18°C.", "blind_spots": "Không biết cá trong tủ đã được đem nấu món gì."}
    ]

    relationships = [
        {"from_agent_id": "h_dr_nghiem", "to_agent_id": "h_nurse_hoa", "trust": 0.85, "reason": "Cộng sự cấp cứu lâu năm"},
        {"from_agent_id": "h_nurse_hoa", "to_agent_id": "h_dr_nghiem", "trust": 0.90, "reason": "Tin tưởng chuyên môn cấp trên"},
        {"from_agent_id": "h_dr_nghiem", "to_agent_id": "h_director_khanh", "trust": 0.80, "reason": "Báo cáo trực tiếp ban giám đốc"},
        {"from_agent_id": "h_director_khanh", "to_agent_id": "h_dr_nghiem", "trust": 0.85, "reason": "Trưởng khoa tin cậy"},
        {"from_agent_id": "h_chef_tam", "to_agent_id": "h_cleaner_bich", "trust": 0.60, "reason": "Nhân viên cùng bộ phận nhà bếp"},
        {"from_agent_id": "h_cleaner_bich", "to_agent_id": "h_vendor_mai", "trust": 0.78, "reason": "Thường ra quán trà đá tâm sự"},
        {"from_agent_id": "h_vendor_mai", "to_agent_id": "h_blogger_khoa", "trust": 0.70, "reason": "Kể chuyện ly kỳ cho tiktoker"},
        {"from_agent_id": "h_blogger_khoa", "to_agent_id": "h_reporter_nga", "trust": 0.65, "reason": "Cùng săn tin độc quyền"},
        {"from_agent_id": "h_reporter_nga", "to_agent_id": "h_officer_nam", "trust": 0.75, "reason": "Nguồn tin cơ quan chức năng"},
        {"from_agent_id": "h_officer_nam", "to_agent_id": "h_guard_dung", "trust": 0.70, "reason": "Phối hợp giữ trật tự cổng viện"},
        {"from_agent_id": "h_guard_dung", "to_agent_id": "h_driver_long", "trust": 0.80, "reason": "Anh em lái xe và bảo vệ"},
        {"from_agent_id": "h_tech_minh", "to_agent_id": "h_inspector_hung", "trust": 0.88, "reason": "Đồng nghiệp khoa học kiểm nghiệm"},
        {"from_agent_id": "h_inspector_hung", "to_agent_id": "h_tech_minh", "trust": 0.85, "reason": "Tin cậy kết quả xét nghiệm vi sinh"},
        {"from_agent_id": "h_relative_quang", "to_agent_id": "h_vendor_mai", "trust": 0.75, "reason": "Hỏi thăm tình hình bên ngoài"},
        {"from_agent_id": "h_pharmacist_yen", "to_agent_id": "h_dr_nghiem", "trust": 0.82, "reason": "Xác nhận phác đồ điều trị"},
        {"from_agent_id": "h_student_linh", "to_agent_id": "h_patient_tri", "trust": 0.80, "reason": "Học việc từ bác sĩ nội trú"},
        {"from_agent_id": "h_nurse_phuong", "to_agent_id": "h_nurse_hoa", "trust": 0.80, "reason": "Đồng nghiệp điều dưỡng"},
        {"from_agent_id": "h_mortuary_van", "to_agent_id": "h_guard_dung", "trust": 0.70, "reason": "Chào hỏi ca trực hàng ngày"},
        {"from_agent_id": "h_technician_son", "to_agent_id": "h_chef_tam", "trust": 0.65, "reason": "Đã từng cảnh báo về tủ lạnh hỏng"}
    ]

    observations = [
        {"agent_id": "h_dr_nghiem", "claim_id": "claim_poisoning_occurred", "confidence": 0.95, "basis": "direct_observation"},
        {"agent_id": "h_dr_nghiem", "claim_id": "claim_zero_fatalities", "confidence": 0.90, "basis": "direct_observation"},
        {"agent_id": "h_dr_nghiem", "claim_id": "claim_no_virus_leak", "confidence": 0.85, "basis": "clinical_judgment"},
        {"agent_id": "h_nurse_hoa", "claim_id": "claim_poisoning_occurred", "confidence": 0.90, "basis": "direct_observation"},
        {"agent_id": "h_nurse_hoa", "claim_id": "claim_zero_fatalities", "confidence": 0.85, "basis": "direct_observation"},
        {"agent_id": "h_tech_minh", "claim_id": "claim_no_virus_leak", "confidence": 0.98, "basis": "direct_observation"},
        {"agent_id": "h_chef_tam", "claim_id": "claim_canteen_fridge_broken", "confidence": 0.80, "basis": "direct_observation"},
        {"agent_id": "h_cleaner_bich", "claim_id": "claim_canteen_fridge_broken", "confidence": 0.90, "basis": "direct_observation"},
        {"agent_id": "h_technician_son", "claim_id": "claim_canteen_fridge_broken", "confidence": 0.95, "basis": "direct_observation"},
        {"agent_id": "h_mortuary_van", "claim_id": "claim_zero_fatalities", "confidence": 0.95, "basis": "direct_observation"},
        {"agent_id": "h_patient_tri", "claim_id": "claim_poisoning_occurred", "confidence": 0.95, "basis": "personal_experience"},
        {"agent_id": "h_patient_tri", "claim_id": "claim_tuna_histamine_origin", "confidence": 0.80, "basis": "medical_knowledge"},
        {"agent_id": "h_pharmacist_yen", "claim_id": "claim_zero_fatalities", "confidence": 0.85, "basis": "drug_requisition_record"}
    ]

    evidence = [
        {
            "id": "ev_canteen_inspection",
            "name": "Biên bản kiểm tra an toàn căng tin",
            "description": "Biên bản ghi nhận tủ đông số 2 bảo quản cá ngừ bị hỏng rơ-le nhiệt, nhiệt độ thực tế 14.8°C khiến vi khuẩn phân giải sinh lượng lớn histamine.",
            "location_id": "nutrition_dept",
            "reliability": 0.98,
            "source_category": "official_inspection",
            "supports": ["claim_canteen_fridge_broken", "claim_tuna_histamine_origin", "claim_poisoning_occurred"],
            "contradicts": ["claim_lab_virus_leak"],
            "unlock_condition": "Kiểm tra kho dinh dưỡng sau khi nói chuyện với Thanh tra Hùng hoặc Kỹ thuật viên Sơn."
        },
        {
            "id": "ev_er_logbook",
            "name": "Sổ nhật ký tiếp nhận cấp cứu",
            "description": "Danh sách 18 ca cấp cứu: Tất cả chẩn đoán 'Hội chứng ngộ độc thức ăn dạng histamine', sinh hiệu ổn định, 0 trường hợp tử vong hay thở máy.",
            "location_id": "er_triage",
            "reliability": 0.99,
            "source_category": "medical_record",
            "supports": ["claim_poisoning_occurred", "claim_zero_fatalities"],
            "contradicts": ["claim_patient_fatalities", "claim_lab_virus_leak"],
            "unlock_condition": "Khám xét bàn phân loại sau khi trao đổi với Bác sĩ Nghiêm hoặc Y tá Hoa."
        },
        {
            "id": "ev_lab_biosafety_audit",
            "name": "Biên bản an toàn sinh học Lab Vi sinh",
            "description": "Chứng nhận phòng xét nghiệm cấp độ 2: Áp suất âm chuẩn, chỉ lưu giữ các chủng vi khuẩn Gram âm/dương thường quy, không có mẫu virus lây nhiễm nguy hiểm.",
            "location_id": "lab_biosafety",
            "reliability": 0.97,
            "source_category": "scientific_audit",
            "supports": ["claim_no_virus_leak"],
            "contradicts": ["claim_lab_virus_leak"],
            "unlock_condition": "Kiểm tra phòng xét nghiệm sau khi phỏng vấn Kỹ thuật viên Minh."
        },
        {
            "id": "ev_mortuary_clearance",
            "name": "Sổ theo dõi bàn giao nhà tang lễ",
            "description": "Nhật ký nhà xác ngày hôm nay ghi nhận 0 ca tiếp nhận mới từ các khoa phòng; không có thi thể nào được chuyển đi trong ngày.",
            "location_id": "mortuary",
            "reliability": 0.99,
            "source_category": "administrative_record",
            "supports": ["claim_zero_fatalities"],
            "contradicts": ["claim_patient_fatalities", "claim_black_bodybags_moved"],
            "unlock_condition": "Kiểm tra nhà tang lễ sau khi hỏi Bác Vạn."
        }
    ]

    public_events = [
        {"game_second": 120, "event_type": "ambulance_surge", "headline": "Xe cấp cứu đưa 6 bác sĩ nhập viện liên tục", "description": "Nhiều y bác sĩ ca trực trưa có biểu hiện đau quặn bụng và dị ứng da cấp.", "location_id": "er_entrance"},
        {"game_second": 300, "event_type": "gate_restricted", "headline": "Bệnh viện dựng rào chắn phân luồng", "description": "Lực lượng bảo vệ và công an tạm ngừng tiếp nhận khám ngoại trú để tập trung cấp cứu.", "location_id": "hospital_lobby"},
        {"game_second": 480, "event_type": "rumor_virus", "headline": "Blogger Khoa livestream: 'Dịch bệnh lạ ở viện Metro?'", "description": "Video quay cảnh bảo hộ y tế lan truyền chóng mặt với hàng nghìn lượt chia sẻ hoảng loạn.", "location_id": "street_cafe"},
        {"game_second": 720, "event_type": "food_inspection", "headline": "Đoàn thanh tra an toàn thực phẩm có mặt", "description": "Lực lượng chức năng tiến hành niêm phong mẫu thức ăn trưa tại căng tin bệnh viện.", "location_id": "canteen"}
    ]

    truth_form = {
        "event_type": {
            "type": "enum",
            "label": "Bản chất sự cố y tế",
            "answer": "food_poisoning",
            "options": [
                {"value": "food_poisoning", "label": "Ngộ độc thực phẩm tập thể"},
                {"value": "viral_outbreak", "label": "Bùng phát virus nguy hiểm"},
                {"value": "chemical_poisoning", "label": "Đầu độc hóa chất"},
                {"value": "unknown", "label": "Chưa xác định"}
            ]
        },
        "cause_category": {
            "type": "enum",
            "label": "Nguyên nhân gốc rễ",
            "answer": "bacterial_toxin",
            "options": [
                {"value": "bacterial_toxin", "label": "Độc tố histamine do bảo quản cá ngừ sai quy chuẩn"},
                {"value": "lab_leak", "label": "Rò rỉ mầm bệnh từ phòng thí nghiệm vi sinh"},
                {"value": "deliberate_sabotage", "label": "Kẻ xấu cố ý phá hoại nguồn thức ăn"},
                {"value": "unknown", "label": "Không rõ"}
            ]
        },
        "major_explosion": {
            "type": "boolean",
            "label": "Có sự cố an toàn sinh học cấp độ cao / thảm họa lây nhiễm không?",
            "answer": False
        },
        "fatalities": {
            "type": "integer",
            "label": "Số người tử vong thực tế",
            "answer": 0,
            "min": 0,
            "max": 50
        },
        "management_issue": {
            "type": "enum",
            "label": "Trách nhiệm quản lý",
            "answer": "broken_equipment_neglected",
            "options": [
                {"value": "broken_equipment_neglected", "label": "Không sửa chữa tủ đông bảo quản thực phẩm bị hỏng"},
                {"value": "illegal_experiment", "label": "Tiến hành thí nghiệm virus trái phép"},
                {"value": "none", "label": "Không có sai phạm"}
            ]
        }
    }

    save_json(s_dir, "scenario.json", metadata)
    save_json(s_dir, "locations.json", locations)
    save_json(s_dir, "claims.json", claims)
    save_json(s_dir, "agents.json", agents)
    save_json(s_dir, "relationships.json", relationships)
    save_json(s_dir, "observations.json", observations)
    save_json(s_dir, "evidence.json", evidence)
    save_json(s_dir, "public-events.json", public_events)
    save_json(s_dir, "truth-form.json", truth_form)

# ==============================================================================
# SCENARIO 3: midtown-bank-run (Cơn Hoảng Loạn Ngân Hàng Midtown)
# ==============================================================================
def create_bank_scenario():
    sc_id = "midtown-bank-run"
    s_dir = os.path.join(BASE_DIR, sc_id)

    metadata = {
        "id": sc_id,
        "title": "Cơn Hoảng Loạn Ngân Hàng Midtown",
        "version": 1,
        "description": "Lúc 14:15, hệ thống máy chủ cơ sở dữ liệu Core Banking Midtown gặp lỗi deadlock kết nối đường truyền, khiến hệ thống ATM và Mobile Banking bị gián đoạn 45 phút. Không mất mát tài sản. Tin đồn thất thiệt lan truyền rằng ngân hàng vỡ nợ, giám đốc đã ôm vàng bỏ trốn và chuẩn bị phong tỏa tài khoản người dân.",
        "primary_false_narrative": {
            "adoption_threshold": 0.65,
            "defeat_ratio": 0.75,
            "claims": [
                {"claim_id": "claim_bank_insolvency", "weight": 0.35},
                {"claim_id": "claim_ceo_fled_with_gold", "weight": 0.40},
                {"claim_id": "claim_account_freeze_order", "weight": 0.25}
            ]
        },
        "required_truth_fields": ["event_type", "cause_category", "major_explosion", "fatalities"],
        "bonus_truth_fields": ["system_status", "management_issue"]
    }

    locations = [
        {"id": "bank_hall", "name": "Đại Sảnh Giao Dịch", "description": "Hàng trăm khách hàng chen lấn xếp hàng rút tiền mặt.", "position": {"x": -24, "y": 0, "z": -16}, "connected_locations": ["atm_gallery", "teller_counters", "vip_lounge"]},
        {"id": "atm_gallery", "name": "Khu Cây Rút Tiền ATM", "description": "Dãy cây ATM màn hình báo lỗi bảo trì hệ thống kết nối.", "position": {"x": -32, "y": 0, "z": -24}, "connected_locations": ["bank_hall", "financial_square"]},
        {"id": "teller_counters", "name": "Quầy Giao Dịch Viên", "description": "Nơi các nhân viên giao dịch đang trấn an khách hàng.", "position": {"x": -20, "y": 0, "z": -26}, "connected_locations": ["bank_hall", "cash_vault"]},
        {"id": "cash_vault", "name": "Hầm Kho Tiền Trung Tâm", "description": "Hầm chứa tiền và vàng dự trữ của ngân hàng được bảo vệ nghiêm ngặt.", "position": {"x": -28, "y": 0, "z": -8}, "connected_locations": ["teller_counters", "security_control"]},
        {"id": "security_control", "name": "Phòng An Ninh Giám Sát", "description": "Phòng điều khiển camera CCTV toàn bộ chi nhánh.", "position": {"x": -14, "y": 0, "z": -18}, "connected_locations": ["cash_vault", "bank_hall"]},
        {"id": "it_server_room", "name": "Trung Tâm Máy Chủ IT", "description": "Phòng máy chủ nơi các kỹ sư khắc phục sự cố phần mềm.", "position": {"x": 2, "y": 0, "z": 6}, "connected_locations": ["bank_hall", "boardroom"]},
        {"id": "boardroom", "name": "Phòng Họp Ban Giám Đốc", "description": "Ban điều hành đang họp khẩn cấp để đưa ra thông cáo báo chí.", "position": {"x": 26, "y": 0, "z": -20}, "connected_locations": ["it_server_room", "vip_lounge"]},
        {"id": "vip_lounge", "name": "Phòng Tiếp Khách VIP", "description": "Các nhà đầu tư lớn đang yêu cầu gặp trực tiếp Tổng Giám đốc.", "position": {"x": -6, "y": 0, "z": -6}, "connected_locations": ["bank_hall", "boardroom"]},
        {"id": "financial_square", "name": "Quảng Trường Tài Chính", "description": "Đông đảo người dân và đám đông hiếu kỳ tụ tập trước cổng ngân hàng.", "position": {"x": 10, "y": 0, "z": -4}, "connected_locations": ["atm_gallery", "broker_cafe"]},
        {"id": "broker_cafe", "name": "Cà Phê Môi Giới Chứng Khoán", "description": "Nơi các tay săn tin đồn và môi giới tài chính bàn tán rôm rả.", "position": {"x": 18, "y": 0, "z": 12}, "connected_locations": ["financial_square", "news_bureau"]},
        {"id": "news_bureau", "name": "Tòa Soạn Bản Tin Thị Trường", "description": "Phóng viên tài chính đang chuẩn bị phát sóng bản tin trực tiếp.", "position": {"x": -22, "y": 0, "z": 14}, "connected_locations": ["broker_cafe"]},
        {"id": "central_bank_desk", "name": "Văn Phòng Đại Diện Ngân Hàng Nhà Nước", "description": "Tổ giám sát của Ngân hàng Trung ương theo dõi thanh khoản.", "position": {"x": -10, "y": 0, "z": 18}, "connected_locations": ["financial_square"]}
    ]

    claims = [
        {"id": "claim_tech_outage_occurred", "text": "Hệ thống kết nối Core Banking bị gián đoạn kỹ thuật tạm thời", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_database_deadlock_cause", "text": "Sự cố bắt nguồn từ xung đột deadlock trong tiến trình nâng cấp bản vá định kỳ", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_liquidity_sufficient", "text": "Hầm tiền vẫn đầy đủ dự trữ thanh khoản, không mất mát tiền gửi", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_ceo_present", "text": "Tổng giám đốc đang có mặt tại phòng họp tầng 3 điều hành khắc phục", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_zero_fatalities", "text": "Không có ai tử vong hay bạo lực xảy ra", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_bank_insolvency", "text": "Ngân hàng Midtown mất khả năng thanh khoản và sắp công bố phá sản", "classification": "false", "role": "primary_false", "weight": 0.35},
        {"id": "claim_ceo_fled_with_gold", "text": "Tổng giám đốc đã bí mật vận chuyển toàn bộ vàng dự trữ và bay sang nước ngoài", "classification": "false", "role": "primary_false", "weight": 0.40},
        {"id": "claim_account_freeze_order", "text": "Tất cả tài khoản tiết kiệm của người dân sẽ bị đóng băng vĩnh viễn từ 17:00", "classification": "false", "role": "primary_false", "weight": 0.25},
        {"id": "claim_secret_police_raid", "text": "Cảnh sát kinh tế đã bao vây bắt giữ toàn bộ ban quản trị", "classification": "false", "role": "secondary_false", "weight": 0.0},
        {"id": "claim_atm_cash_empty_forever", "text": "ATM sẽ không bao giờ nhả tiền nữa vì máy chủ đã bị xóa sạch dữ liệu", "classification": "false", "role": "secondary_false", "weight": 0.0}
    ]

    agents = [
        {"id": "b_ceo_viet", "name": "Tổng Giám đốc Viết", "role": "Tổng Giám đốc Ngân hàng Midtown", "location_id": "boardroom", "influence": 0.85, "credibility": 0.78, "skepticism": 0.80, "sociality": 0.35, "deception_tendency": 0.10, "receptiveness": 0.30, "personality": "quyết liệt, tự tin, chịu áp lực cao", "private_motive": "Bảo vệ uy tín thương hiệu và ngăn chặn làn sóng rút tiền ồ ạt.", "direct_knowledge": "Ngân hàng có thặng dư tiền mặt 2.000 tỷ VNĐ trong kho, sự cố chỉ là do nghẽn mạch IT.", "blind_spots": "Chưa nắm được mức độ hoang mang thực tế ở ngoài sảnh."},
        {"id": "b_cio_thang", "name": "Giám đốc IT Thắng", "role": "Giám đốc Công nghệ Thông tin", "location_id": "it_server_room", "influence": 0.60, "credibility": 0.92, "skepticism": 0.85, "sociality": 0.25, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "khoa học, logic, thức trắng đêm", "private_motive": "Rollback bản vá lỗi và khởi động lại dịch vụ nhanh nhất.", "direct_knowledge": "Lỗi do bản cập nhật cơ sở dữ liệu phiên bản 4.12 gây xung đột lock table, dữ liệu số dư 100% an toàn.", "blind_spots": "Không giao tiếp trực tiếp với khách hàng."},
        {"id": "b_vault_chief_hung", "name": "Trưởng kho Hùng", "role": "Trưởng ban Quản lý Kho quỹ", "location_id": "cash_vault", "influence": 0.45, "credibility": 0.88, "skepticism": 0.75, "sociality": 0.30, "deception_tendency": 0.05, "receptiveness": 0.40, "personality": "nguyên tắc, cẩn mật, chuẩn chỉ", "private_motive": "Bảo vệ an toàn tuyệt đối cho hầm vàng và kho tiền mặt.", "direct_knowledge": "Kho quỹ đầy đủ 100% tài sản theo sổ sách kiểm kê lúc 14:00, không có xe chở tiền nào rời kho.", "blind_spots": "Không biết hệ thống mạng Internet bên ngoài."},
        {"id": "b_teller_mai", "name": "Giao dịch viên Mai", "role": "Kiểm soát viên sảnh giao dịch", "location_id": "teller_counters", "influence": 0.40, "credibility": 0.70, "skepticism": 0.45, "sociality": 0.80, "deception_tendency": 0.08, "receptiveness": 0.65, "personality": "tận tụy, kiên nhẫn, mệt mỏi", "private_motive": "Xoa dịu khách hàng cao tuổi đến rút tiền.", "direct_knowledge": "Các lệnh giao dịch tại quầy vẫn ký tay và đối soát bình thường khi mạng phục hồi.", "blind_spots": "Chỉ nhìn thấy màn hình máy tính báo vòng xoay chờ phản hồi."},
        {"id": "b_investor_lam", "name": "Doanh nhân Lâm", "role": "Khách hàng VIP / Chủ doanh nghiệp", "location_id": "vip_lounge", "influence": 0.65, "credibility": 0.60, "skepticism": 0.35, "sociality": 0.75, "deception_tendency": 0.15, "receptiveness": 0.80, "personality": "nóng nảy, đa nghi, quyền lực", "private_motive": "Muốn chuyển gấp 50 tỷ tiền gửi sang ngân hàng nước ngoài.", "direct_knowledge": "Ứng dụng VIP Banking báo 'Không thể kết nối máy chủ' lúc 14:15.", "blind_spots": "Tin vào tin đồn trên nhóm chat ngầm rằng ngân hàng sắp sập."},
        {"id": "b_security_guard", "name": "Bảo vệ Tuấn", "role": "Đội trưởng An ninh", "location_id": "bank_hall", "influence": 0.35, "credibility": 0.65, "skepticism": 0.55, "sociality": 0.70, "deception_tendency": 0.05, "receptiveness": 0.60, "personality": "nghiêm nghị, cảnh giác", "private_motive": "Ngăn chặn cảnh chen lấn xô đẩy phá vỡ cửa kính ngân hàng.", "direct_knowledge": "Khách hàng dồn về sảnh quá đông khiến cửa tự động phải chuyển sang mở thủ công.", "blind_spots": "Không rõ nguyên nhân kỹ thuật bên trong."},
        {"id": "b_broker_trung", "name": "Môi giới Trung", "role": "Chuyên viên Chứng khoán", "location_id": "broker_cafe", "influence": 0.70, "credibility": 0.45, "skepticism": 0.30, "sociality": 0.90, "deception_tendency": 0.30, "receptiveness": 0.75, "personality": "thích đầu cơ, lan truyền tin đồn thị trường", "private_motive": "Tạo sóng hoảng loạn để gom cổ phiếu ngân hàng giá rẻ.", "direct_knowledge": "Cổ phiếu Midtown giảm sàn 7% chỉ sau 15 phút mở phiên chiều.", "blind_spots": "Không có bằng chứng kiểm toán thực tế."},
        {"id": "b_inspector_hoa", "name": "Thanh tra Hoa", "role": "Giám sát Ngân hàng Nhà nước", "location_id": "central_bank_desk", "influence": 0.75, "credibility": 0.95, "skepticism": 0.90, "sociality": 0.35, "deception_tendency": 0.02, "receptiveness": 0.30, "personality": "nghiêm túc, thận trọng, uy tín", "private_motive": "Giữ vững an toàn hệ thống tài chính quốc gia.", "direct_knowledge": "Chỉ số an toàn vốn CAR của Midtown đạt 13.5%, cao hơn mức quy định; NHNN sẵn sàng cấp thanh khoản nếu cần.", "blind_spots": "Chậm trễ trong việc ban hành văn bản đính chính chính thức."},
        {"id": "b_citizen_ba", "name": "Bà Ba Hưu Trí", "role": "Người gửi tiền tiết kiệm", "location_id": "atm_gallery", "influence": 0.30, "credibility": 0.50, "skepticism": 0.20, "sociality": 0.85, "deception_tendency": 0.05, "receptiveness": 0.90, "personality": "lo lắng, dễ xúc động, truyền tai nhanh", "private_motive": "Bảo vệ toàn bộ số tiền dưỡng già tích cóp cả đời.", "direct_knowledge": "Đút thẻ vào 3 cây ATM đều bị nhả ra và báo lỗi liên lạc.", "blind_spots": "Nghe người đứng cạnh bảo sếp ngân hàng trốn mất tiêu rồi."},
        {"id": "b_reporter_tuan", "name": "Nhà báo Tuấn", "role": "Phóng viên Kinh tế - Tài chính", "location_id": "news_bureau", "influence": 0.68, "credibility": 0.75, "skepticism": 0.65, "sociality": 0.70, "deception_tendency": 0.08, "receptiveness": 0.60, "personality": "nhanh nhạy, khách quan, kiểm chứng", "private_motive": "Đưa tin xác thực để bảo vệ thị trường.", "direct_knowledge": "Đang liên hệ phòng truyền thông ngân hàng để lấy thông cáo chính thức.", "blind_spots": "Bị áp lực bởi các trang tin lá cải đăng bài câu view giật gân."},
        {"id": "b_it_engineer_duc", "name": "Kỹ sư Đức", "role": "Chuyên viên Quản trị Cơ sở dữ liệu", "location_id": "it_server_room", "influence": 0.40, "credibility": 0.85, "skepticism": 0.80, "sociality": 0.30, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "thực tế, ít nói, làm việc chính xác", "private_motive": "Khôi phục trạng thái database cluster.", "direct_knowledge": "Quá trình rollback đã hoàn tất 80%, các node đang đồng bộ lại dữ liệu nhật ký giao dịch.", "blind_spots": "Không biết tin đồn bên ngoài đã lan tới đâu."},
        {"id": "b_lawyer_phuc", "name": "Luật sư Phúc", "role": "Cố vấn Pháp chế Ngân hàng", "location_id": "boardroom", "influence": 0.55, "credibility": 0.82, "skepticism": 0.75, "sociality": 0.45, "deception_tendency": 0.05, "receptiveness": 0.40, "personality": "chắc chắn, tuân thủ pháp luật", "private_motive": "Soạn thảo văn bản cảnh cáo các đối tượng tung tin bịa đặt.", "direct_knowledge": "Khẳng định không có bất kỳ lệnh phong tỏa hay khởi tố nào từ cơ quan điều tra.", "blind_spots": "Không am hiểu sâu về kỹ thuật máy chủ."},
        {"id": "b_cctv_operator_hieu", "name": "Nhân viên Hiếu", "role": "Trực ban Camera An ninh", "location_id": "security_control", "influence": 0.35, "credibility": 0.75, "skepticism": 0.65, "sociality": 0.50, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "quan sát tinh tường", "private_motive": "Giám sát an ninh toàn bộ các cửa ra vào.", "direct_knowledge": "Camera hầm tiền cho thấy cửa hầm khóa nguyên tem niêm phong, Tổng Giám đốc đi thang máy lên tầng 3 chứ không đi đâu.", "blind_spots": "Chỉ quan sát trong tòa nhà."},
        {"id": "b_vendor_dung", "name": "Chú Dũng Xe Ôm", "role": "Tài xế trước cổng ngân hàng", "location_id": "financial_square", "influence": 0.42, "credibility": 0.40, "skepticism": 0.25, "sociality": 0.95, "deception_tendency": 0.20, "receptiveness": 0.85, "personality": "thích hóng chuyện, khuếch đại thông tin", "private_motive": "Chở khách đi rút tiền ở các chi nhánh khác để kiếm cuốc xe.", "direct_knowledge": "Thấy nhiều người kéo đến cổng khóc lóc và đập cửa kính.", "blind_spots": "Tự tưởng tượng ra cảnh ngân hàng vỡ nợ như trong phim."},
        {"id": "b_customer_lan", "name": "Chị Lan Kế Toán", "role": "Kế toán trưởng công ty tư nhân", "location_id": "bank_hall", "influence": 0.48, "credibility": 0.65, "skepticism": 0.50, "sociality": 0.70, "deception_tendency": 0.05, "receptiveness": 0.70, "personality": "sốt ruột, bồn chồn", "private_motive": "Phải chuyển tiền trả lương cho 200 công nhân trước 16:30.", "direct_knowledge": "Lệnh ủy nhiệm chi nộp tại quầy đã được nhân viên tiếp nhận và hứa duyệt ngay khi máy chủ thông mạng.", "blind_spots": "Sợ bị công nhân đình công nếu tiền không về tài khoản."},
        {"id": "b_banker_phong", "name": "Phó Giám đốc Phong", "role": "Phó Tổng Giám đốc Khối Vận hành", "location_id": "boardroom", "influence": 0.70, "credibility": 0.78, "skepticism": 0.70, "sociality": 0.50, "deception_tendency": 0.08, "receptiveness": 0.45, "personality": "thực chiến, linh hoạt", "private_motive": "Điều phối xe chở thêm 300 tỷ tiền mặt từ kho dự trữ quốc gia về chi nhánh.", "direct_knowledge": "Đã có lệnh điều xe tiếp quỹ tiền mặt để đáp ứng nhu cầu rút tiền tại quầy.", "blind_spots": "Chưa thông báo kịp ra loa phát thanh sảnh."},
        {"id": "b_student_khoa", "name": "Sinh viên Khoa", "role": "Sinh viên Kinh tế", "location_id": "broker_cafe", "influence": 0.38, "credibility": 0.55, "skepticism": 0.40, "sociality": 0.80, "deception_tendency": 0.10, "receptiveness": 0.75, "personality": "nhanh nhẹn, lướt mạng xã hội liên tục", "private_motive": "Thực tập nghiên cứu về hành vi tâm lý đám đông bank-run.", "direct_knowledge": "Thấy trên các hội nhóm Facebook bắt đầu lan truyền ảnh chế Tổng giám đốc bỏ trốn.", "blind_spots": "Không có khả năng kiểm chứng nguồn tin nội bộ."},
        {"id": "b_driver_viet", "name": "Lái xe Việt", "role": "Tài xế xe bọc thép chở tiền", "location_id": "cash_vault", "influence": 0.35, "credibility": 0.75, "skepticism": 0.60, "sociality": 0.55, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "trầm tính, bảo mật", "private_motive": "Tuân thủ nghiêm ngặt quy trình an ninh vận chuyển.", "direct_knowledge": "Xe chở tiền vẫn đỗ trong hầm ngầm, sẵn sàng nạp thêm tiền vào các trụ ATM.", "blind_spots": "Không được phép tiết lộ lộ trình xe."},
        {"id": "b_teller_nga", "name": "Giao dịch viên Nga", "role": "Nhân viên chăm sóc khách hàng", "location_id": "vip_lounge", "influence": 0.42, "credibility": 0.72, "skepticism": 0.50, "sociality": 0.75, "deception_tendency": 0.05, "receptiveness": 0.60, "personality": "nhã nhặn, khéo léo", "private_motive": "Giữ chân các khách hàng gửi tiền tỷ không rút trước hạn.", "direct_knowledge": "Cam kết lãi suất tiền gửi vẫn được đảm bảo đầy đủ theo quy định của Ngân hàng Trung ương.", "blind_spots": "Bị khách hàng VIP mắng mỏ gây áp lực tâm lý."},
        {"id": "b_security_minh", "name": "Bảo vệ Minh", "role": "Bảo vệ khu vực ATM", "location_id": "atm_gallery", "influence": 0.32, "credibility": 0.60, "skepticism": 0.40, "sociality": 0.70, "deception_tendency": 0.05, "receptiveness": 0.65, "personality": "nhiệt tình, vất vả", "private_motive": "Hướng dẫn người dân xếp hàng trật tự.", "direct_knowledge": "Kỹ thuật viên IT vừa thông báo qua bộ đàm là máy chủ đang boot lại.", "blind_spots": "Không giải thích được thuật ngữ deadlock cho người dân hiểu."}
    ]

    relationships = [
        {"from_agent_id": "b_ceo_viet", "to_agent_id": "b_cio_thang", "trust": 0.88, "reason": "Tin cậy năng lực giám đốc IT"},
        {"from_agent_id": "b_cio_thang", "to_agent_id": "b_it_engineer_duc", "trust": 0.90, "reason": "Cộng sự trực tiếp phòng máy chủ"},
        {"from_agent_id": "b_ceo_viet", "to_agent_id": "b_vault_chief_hung", "trust": 0.92, "reason": "Bảo chứng tài sản kho quỹ"},
        {"from_agent_id": "b_inspector_hoa", "to_agent_id": "b_ceo_viet", "trust": 0.80, "reason": "Quan hệ giám sát quản lý nhà nước"},
        {"from_agent_id": "b_cctv_operator_hieu", "to_agent_id": "b_security_guard", "trust": 0.85, "reason": "Phối hợp an ninh nội bộ"},
        {"from_agent_id": "b_teller_mai", "to_agent_id": "b_banker_phong", "trust": 0.80, "reason": "Nhận chỉ đạo điều phối tiền mặt"},
        {"from_agent_id": "b_vendor_dung", "to_agent_id": "b_broker_trung", "trust": 0.70, "reason": "Hóng chuyện thị trường ở quán cafe"},
        {"from_agent_id": "b_broker_trung", "to_agent_id": "b_student_khoa", "trust": 0.65, "reason": "Bàn luận phân tích chứng khoán"},
        {"from_agent_id": "b_investor_lam", "to_agent_id": "b_ceo_viet", "trust": 0.75, "reason": "Mối quan hệ khách hàng chiến lược"},
        {"from_agent_id": "b_reporter_tuan", "to_agent_id": "b_inspector_hoa", "trust": 0.85, "reason": "Nguồn tin cơ quan thanh tra"},
        {"from_agent_id": "b_citizen_ba", "to_agent_id": "b_vendor_dung", "trust": 0.75, "reason": "Người quen cùng khu phố"},
        {"from_agent_id": "b_customer_lan", "to_agent_id": "b_teller_mai", "trust": 0.78, "reason": "Làm việc thanh toán thường xuyên"},
        {"from_agent_id": "b_driver_viet", "to_agent_id": "b_vault_chief_hung", "trust": 0.90, "reason": "Đồng nghiệp áp tải tiền mặt"},
        {"from_agent_id": "b_lawyer_phuc", "to_agent_id": "b_ceo_viet", "trust": 0.85, "reason": "Cố vấn pháp chế trưởng"}
    ]

    observations = [
        {"agent_id": "b_ceo_viet", "claim_id": "claim_ceo_present", "confidence": 0.99, "basis": "direct_observation"},
        {"agent_id": "b_ceo_viet", "claim_id": "claim_liquidity_sufficient", "confidence": 0.95, "basis": "financial_records"},
        {"agent_id": "b_cio_thang", "claim_id": "claim_tech_outage_occurred", "confidence": 0.99, "basis": "server_logs"},
        {"agent_id": "b_cio_thang", "claim_id": "claim_database_deadlock_cause", "confidence": 0.95, "basis": "direct_observation"},
        {"agent_id": "b_vault_chief_hung", "claim_id": "claim_liquidity_sufficient", "confidence": 0.98, "basis": "direct_observation"},
        {"agent_id": "b_cctv_operator_hieu", "claim_id": "claim_ceo_present", "confidence": 0.95, "basis": "cctv_monitoring"},
        {"agent_id": "b_inspector_hoa", "claim_id": "claim_liquidity_sufficient", "confidence": 0.92, "basis": "supervisory_data"},
        {"agent_id": "b_it_engineer_duc", "claim_id": "claim_database_deadlock_cause", "confidence": 0.98, "basis": "debug_console"}
    ]

    evidence = [
        {
            "id": "ev_it_incident_log",
            "name": "Nhật ký hệ thống máy chủ IT (Server Syslog)",
            "description": "Log máy chủ Core Banking ghi nhận: 'Deadlock on transaction table during patch v4.12 migration, auto-rollback initiated at 14:18'. Không có xâm nhập dữ liệu trái phép.",
            "location_id": "it_server_room",
            "reliability": 0.99,
            "source_category": "technical_log",
            "supports": ["claim_tech_outage_occurred", "claim_database_deadlock_cause"],
            "contradicts": ["claim_bank_insolvency", "claim_atm_cash_empty_forever"],
            "unlock_condition": "Kiểm tra phòng máy chủ sau khi trao đổi với Giám đốc IT Thắng hoặc Kỹ sư Đức."
        },
        {
            "id": "ev_central_vault_audit",
            "name": "Biên bản kiểm kê kho quỹ khẩn cấp",
            "description": "Biên bản kiểm quỹ có chữ ký liên ngành lúc 14:30: Hầm tiền Midtown lưu trữ 2.150 tỷ VNĐ tiền mặt và 1.200 cây vàng SJC nguyên đai nguyên kiện, đủ đáp ứng 100% thanh khoản.",
            "location_id": "cash_vault",
            "reliability": 0.99,
            "source_category": "official_audit",
            "supports": ["claim_liquidity_sufficient"],
            "contradicts": ["claim_bank_insolvency", "claim_ceo_fled_with_gold"],
            "unlock_condition": "Xem biên bản trong hầm tiền sau khi nói chuyện với Trưởng kho Hùng."
        },
        {
            "id": "ev_central_bank_bulletin",
            "name": "Công điện khẩn Ngân hàng Nhà nước",
            "description": "Công điện xác nhận Ngân hàng Midtown hoạt động an toàn, sự cố kỹ thuật ATM đang được khắc phục và khẳng định quyền lợi người gửi tiền được Nhà nước bảo đảm tuyệt đối.",
            "location_id": "central_bank_desk",
            "reliability": 0.98,
            "source_category": "regulatory_statement",
            "supports": ["claim_liquidity_sufficient", "claim_tech_outage_occurred"],
            "contradicts": ["claim_account_freeze_order", "claim_bank_insolvency"],
            "unlock_condition": "Lấy tại bàn giám sát sau khi phỏng vấn Thanh tra Hoa."
        },
        {
            "id": "ev_cctv_ceo_presence",
            "name": "Băng ghi hình CCTV phòng họp ban lãnh đạo",
            "description": "Camera an ninh xác thực Tổng Giám đốc Viết liên tục có mặt tại phòng họp tầng 3 chỉ đạo ứng phó từ 14:00 đến hiện tại, hoàn toàn không rời khỏi tòa nhà.",
            "location_id": "security_control",
            "reliability": 0.96,
            "source_category": "visual_record",
            "supports": ["claim_ceo_present"],
            "contradicts": ["claim_ceo_fled_with_gold"],
            "unlock_condition": "Trích xuất tại phòng an ninh sau khi trao đổi với Nhân viên Hiếu."
        }
    ]

    public_events = [
        {"game_second": 120, "event_type": "atm_glitch", "headline": "Các cây ATM Midtown đồng loạt báo lỗi kết nối", "description": "Hàng trăm khách hàng không thể rút tiền mặt hoặc kiểm tra số dư qua ứng dụng di động.", "location_id": "atm_gallery"},
        {"game_second": 300, "event_type": "crowd_gathering", "headline": "Đám đông kéo đến đại sảnh ngân hàng", "description": "Người dân lo lắng tập trung yêu cầu rút tiền tiết kiệm trước hạn.", "location_id": "bank_hall"},
        {"game_second": 480, "event_type": "viral_rumor_insolvency", "headline": "Tin đồn 'Ngân hàng phá sản, sếp ôm vàng tháo chạy'", "description": "Nhiều hội nhóm mạng xã hội lan truyền thông tin thất thiệt chưa được kiểm chứng.", "location_id": "broker_cafe"},
        {"game_second": 720, "event_type": "state_bank_announcement", "headline": "Ngân hàng Nhà nước phát đi thông điệp trấn an", "description": "Cơ quan quản lý khẳng định hệ thống Midtown an toàn và thanh khoản dồi dào.", "location_id": "central_bank_desk"}
    ]

    truth_form = {
        "event_type": {
            "type": "enum",
            "label": "Bản chất sự cố ngân hàng",
            "answer": "tech_outage",
            "options": [
                {"value": "tech_outage", "label": "Gián đoạn kỹ thuật hệ thống Core Banking"},
                {"value": "bank_insolvency", "label": "Ngân hàng mất thanh khoản / phá sản"},
                {"value": "cyber_heist", "label": "Vụ trộm mạng quy mô lớn"},
                {"value": "unknown", "label": "Chưa xác định"}
            ]
        },
        "cause_category": {
            "type": "enum",
            "label": "Nguyên nhân chính",
            "answer": "core_banking_deadlock",
            "options": [
                {"value": "core_banking_deadlock", "label": "Xung đột deadlock cơ sở dữ liệu khi nâng cấp bản vá v4.12"},
                {"value": "embezzlement", "label": "Ban lãnh đạo tẩu tán tài sản"},
                {"value": "hacker_ransomware", "label": "Mã độc tống tiền khóa dữ liệu"},
                {"value": "unknown", "label": "Không rõ"}
            ]
        },
        "major_explosion": {
            "type": "boolean",
            "label": "Có sự cố mất trắng tiền gửi / kho quỹ bị cướp phá không?",
            "answer": False
        },
        "fatalities": {
            "type": "integer",
            "label": "Số người thương vong thực tế",
            "answer": 0,
            "min": 0,
            "max": 10
        },
        "management_issue": {
            "type": "enum",
            "label": "Trách nhiệm quản trị",
            "answer": "untested_migration_in_production",
            "options": [
                {"value": "untested_migration_in_production", "label": "Nâng cấp hệ thống trực tiếp trong giờ giao dịch mà chưa test tải kỹ"},
                {"value": "fraudulent_operations", "label": "Hành vi lừa đảo chiếm đoạt tài sản"},
                {"value": "none", "label": "Không có sai sót"}
            ]
        }
    }

    save_json(s_dir, "scenario.json", metadata)
    save_json(s_dir, "locations.json", locations)
    save_json(s_dir, "claims.json", claims)
    save_json(s_dir, "agents.json", agents)
    save_json(s_dir, "relationships.json", relationships)
    save_json(s_dir, "observations.json", observations)
    save_json(s_dir, "evidence.json", evidence)
    save_json(s_dir, "public-events.json", public_events)
    save_json(s_dir, "truth-form.json", truth_form)

# ==============================================================================
# SCENARIO 4: subway-line3-standstill (Chuyến Tàu Ngầm Tuyến Số 3)
# ==============================================================================
def create_subway_scenario():
    sc_id = "subway-line3-standstill"
    s_dir = os.path.join(BASE_DIR, sc_id)

    metadata = {
        "id": sc_id,
        "title": "Chuyến Tàu Ngầm Tuyến Số 3",
        "version": 1,
        "description": "Lúc 18:30 giờ cao điểm, chập cáp tín hiệu đường ray hầm ngầm khiến đoàn tàu Metro E304 kích hoạt hệ thống phanh hãm an toàn khẩn cấp, dừng lại giữa đường hầm. Hệ thống quạt gió khẩn cấp hoạt động tốt, 120 hành khách được sơ tán an toàn theo lối thoát hiểm. 0 ca tử vong. Tin đồn lan truyền rằng tàu bị khủng bố bom khí độc và hàng trăm người ngạt thở kẹt dưới hầm.",
        "primary_false_narrative": {
            "adoption_threshold": 0.65,
            "defeat_ratio": 0.75,
            "claims": [
                {"claim_id": "claim_terrorist_gas_attack", "weight": 0.35},
                {"claim_id": "claim_hundreds_suffocated", "weight": 0.40},
                {"claim_id": "claim_tunnel_collapse_denial", "weight": 0.25}
            ]
        },
        "required_truth_fields": ["event_type", "cause_category", "major_explosion", "fatalities"],
        "bonus_truth_fields": ["evacuation_status", "management_issue"]
    }

    locations = [
        {"id": "station_concourse", "name": "Sảnh Ga Trung Tâm", "description": "Nơi hành khách và thân nhân chen chúc chờ thông tin.", "position": {"x": -24, "y": 0, "z": -16}, "connected_locations": ["ticket_gates", "platform_1", "security_post"]},
        {"id": "ticket_gates", "name": "Cổng Soát Vé Ga", "description": "Cổng an ninh đã mở tự do để sơ tán đám đông.", "position": {"x": -32, "y": 0, "z": -24}, "connected_locations": ["station_concourse", "street_entrance"]},
        {"id": "platform_1", "name": "Sân Ga Ngầm Số 1", "description": "Sân ga dẫn vào hầm tối nơi đoàn tàu đang dừng lại.", "position": {"x": -20, "y": 0, "z": -26}, "connected_locations": ["station_concourse", "train_tunnel_section", "emergency_exit"]},
        {"id": "train_tunnel_section", "name": "Đoạn Hầm Tàu Dừng (Khoang E304)", "description": "Hiện trường đoàn tàu dừng bánh, đèn chiếu sáng khẩn cấp đang bật.", "position": {"x": -28, "y": 0, "z": -8}, "connected_locations": ["platform_1", "track_signal_box"]},
        {"id": "track_signal_box", "name": "Tủ Tín Hiệu Đường Ray Hầm", "description": "Vị trí hộp cáp tín hiệu điều khiển cảm biến an toàn bị chập cháy.", "position": {"x": -14, "y": 0, "z": -18}, "connected_locations": ["train_tunnel_section"]},
        {"id": "station_control_room", "name": "Trung Tâm Điều Hành Ga (OCC)", "description": "Màn hình giám sát tín hiệu toàn tuyến đường sắt đô thị.", "position": {"x": 2, "y": 0, "z": 6}, "connected_locations": ["station_concourse", "ventilation_plant"]},
        {"id": "ventilation_plant", "name": "Trạm Quạt Gió Khẩn Cấp", "description": "Hệ thống quạt hút áp lực đẩy không khí tươi vào đường hầm.", "position": {"x": 26, "y": 0, "z": -20}, "connected_locations": ["station_control_room"]},
        {"id": "emergency_exit", "name": "Lối Thoát Hiểm Số 4", "description": "Cầu thang thoát hiểm phụ dẫn hành khách lên mặt đất.", "position": {"x": -6, "y": 0, "z": -6}, "connected_locations": ["platform_1", "street_entrance"]},
        {"id": "street_entrance", "name": "Cửa Ga Mặt Đường", "description": "Nơi tập trung xe cứu hỏa, cứu thương và cảnh sát phong tỏa.", "position": {"x": 10, "y": 0, "z": -4}, "connected_locations": ["ticket_gates", "emergency_exit", "subway_kiosk"]},
        {"id": "subway_kiosk", "name": "Quầy Báo Cửa Ga", "description": "Nơi người dân tụ tập nghe ngóng và bàn tán tin sốt dẻo.", "position": {"x": 18, "y": 0, "z": 12}, "connected_locations": ["street_entrance", "field_triage_tent"]},
        {"id": "field_triage_tent", "name": "Lều Sơ Cứu Dã Chiến", "description": "Nhân viên y tế 115 tiếp nhận sơ cứu hành khách bị hoảng loạn.", "position": {"x": -22, "y": 0, "z": 14}, "connected_locations": ["street_entrance"]},
        {"id": "security_post", "name": "Chốt Cảnh Sát Ga Ngầm", "description": "Lực lượng cảnh sát bảo vệ an ninh trật tự công cộng ga.", "position": {"x": -10, "y": 0, "z": 18}, "connected_locations": ["station_concourse"]}
    ]

    claims = [
        {"id": "claim_train_stopped_safely", "text": "Đoàn tàu E304 tự động phanh dừng khẩn cấp an toàn trong hầm", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_signal_cable_short", "text": "Nguyên nhân do chập đoản mạch cáp tín hiệu hộp chuyển hướng ray", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_ventilation_working", "text": "Hệ thống thông gió hầm hoạt động hết công suất, cung cấp đủ dưỡng khí", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_all_evacuated_zero_deaths", "text": "120 hành khách được sơ tán an toàn, 0 người tử vong", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_no_structural_collapse", "text": "Kết cấu đường hầm bê tông kiên cố, hoàn toàn không bị sập", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_terrorist_gas_attack", "text": "Có vụ tấn công khủng bố bằng thiết bị nổ khí độc thần kinh trên tàu", "classification": "false", "role": "primary_false", "weight": 0.35},
        {"id": "claim_hundreds_suffocated", "text": "Hàng trăm hành khách đã ngạt thở tử vong kẹt cứng trong các toa tàu kín", "classification": "false", "role": "primary_false", "weight": 0.40},
        {"id": "claim_tunnel_collapse_denial", "text": "Đường hầm đã sập một đoạn và nhà quản lý đường sắt đang chối bỏ", "classification": "false", "role": "primary_false", "weight": 0.25},
        {"id": "claim_smoke_is_chemical_weapon", "text": "Mùi khét bốc lên là vũ khí hóa học hủy diệt", "classification": "false", "role": "secondary_false", "weight": 0.0},
        {"id": "claim_emergency_doors_welded", "text": "Cửa thoát hiểm bị khóa chặt không ai ra được", "classification": "false", "role": "secondary_false", "weight": 0.0}
    ]

    agents = [
        {"id": "s_driver_hai", "name": "Lái tàu Hải", "role": "Trưởng lái tàu E304", "location_id": "train_tunnel_section", "influence": 0.65, "credibility": 0.88, "skepticism": 0.75, "sociality": 0.40, "deception_tendency": 0.05, "receptiveness": 0.45, "personality": "bình tĩnh, giàu kinh nghiệm, trách nhiệm", "private_motive": "Bảo vệ an toàn cho toàn bộ hành khách trên tàu.", "direct_knowledge": "Tàu phanh gấp do mất tín hiệu đèn xanh tự động, trên tàu không có khói độc hay cháy nổ.", "blind_spots": "Không biết tình hình náo loạn trên mặt đất."},
        {"id": "s_dispatcher_trang", "name": "Điều độ viên Trang", "role": "Trưởng ca Trung tâm Điều hành (OCC)", "location_id": "station_control_room", "influence": 0.75, "credibility": 0.90, "skepticism": 0.85, "sociality": 0.35, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "quyết đoán, kỷ luật cao", "private_motive": "Kích hoạt đúng quy trình sơ tán sự cố khẩn cấp.", "direct_knowledge": "Đã bật quạt thông gió hầm chế độ khẩn cấp và cắt điện ray tiếp xúc để đảm bảo an toàn cho hành khách đi bộ.", "blind_spots": "Chưa kiểm tra trực tiếp hiện trường hộp tín hiệu."},
        {"id": "s_tech_quoc", "name": "Kỹ sư Quốc", "role": "Chuyên viên Tín hiệu Đường ray", "location_id": "track_signal_box", "influence": 0.50, "credibility": 0.92, "skepticism": 0.80, "sociality": 0.30, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "thực chứng, cẩn trọng", "private_motive": "Sửa chữa cáp tín hiệu để thông tuyến sớm.", "direct_knowledge": "Hộp cáp tín hiệu số 14 bị chập cháy lớp vỏ cách điện do ẩm mốc, phát ra mùi khét nhựa đặc trưng.", "blind_spots": "Không biết mùi nhựa khét bị đồn thành khí độc."},
        {"id": "s_fire_commander_kien", "name": "Trung tá Kiên", "role": "Chỉ huy Cứu hỏa & Cứu hộ Cứu nạn", "location_id": "street_entrance", "influence": 0.70, "credibility": 0.95, "skepticism": 0.90, "sociality": 0.45, "deception_tendency": 0.02, "receptiveness": 0.30, "personality": "thép, quả cảm, trung thực", "private_motive": "Hoàn thành chiến dịch giải cứu trong thời gian ngắn nhất.", "direct_knowledge": "Lực lượng cứu hộ đã tiếp cận đoàn tàu, dẫn bộ hành khách ra cửa thoát hiểm số 4 an toàn.", "blind_spots": "Chưa nắm hết tâm lý đám đông bên ngoài ga."},
        {"id": "s_passenger_nga", "name": "Hành khách Nga", "role": "Nhân viên văn phòng đi tàu", "location_id": "field_triage_tent", "influence": 0.45, "credibility": 0.65, "skepticism": 0.40, "sociality": 0.80, "deception_tendency": 0.10, "receptiveness": 0.75, "personality": "hoảng loạn, khóc nức nở", "private_motive": "Muốn về nhà với con sau sự cố kinh hoàng.", "direct_knowledge": "Tàu phanh giật cục trong bóng tối, ngửi thấy mùi khét lẹt nên tưởng cháy nổ lớn.", "blind_spots": "Tự suy diễn mùi nhựa khét là hơi độc hóa học."},
        {"id": "s_station_master_vinh", "name": "Trưởng ga Vinh", "role": "Trưởng Ga Trung Tâm", "location_id": "station_concourse", "influence": 0.60, "credibility": 0.80, "skepticism": 0.70, "sociality": 0.55, "deception_tendency": 0.08, "receptiveness": 0.50, "personality": "nỗ lực, căng thẳng", "private_motive": "Mở toàn bộ cửa thoát và hướng dẫn đám đông thoát hiểm.", "direct_knowledge": "Đã mở tất cả cửa soát vé tự do, hệ thống loa phát thanh liên tục trấn an.", "blind_spots": "Loa phát thanh bị tiếng ồn át đi khiến người dân không nghe rõ."},
        {"id": "s_police_son", "name": "Đại úy Sơn", "role": "Cảnh sát Ga ngầm", "location_id": "security_post", "influence": 0.55, "credibility": 0.85, "skepticism": 0.80, "sociality": 0.40, "deception_tendency": 0.02, "receptiveness": 0.40, "personality": "nghiêm nghị, giữ kỷ cương", "private_motive": "Không để xảy ra cảnh giẫm đạp ở cầu thang thoát hiểm.", "direct_knowledge": "Cảnh sát không phát hiện bất kỳ vật thể hay chất nổ khả nghi nào.", "blind_spots": "Không có mặt trong đường hầm sâu."},
        {"id": "s_doctor_huyen", "name": "Bác sĩ Huyền", "role": "Bác sĩ Cấp cứu 115", "location_id": "field_triage_tent", "influence": 0.60, "credibility": 0.90, "skepticism": 0.85, "sociality": 0.45, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "chu đáo, nhân ái, chuyên môn vững", "private_motive": "Sơ cứu cho các ca tụt huyết áp và hoảng sợ.", "direct_knowledge": "Tiếp nhận 4 ca trầy xước nhẹ và khó thở do hoảng loạn tâm lý, 0 ca nhiễm độc hô hấp, 0 người chết.", "blind_spots": "Không tiếp cận hiện trường dưới hầm."},
        {"id": "s_kiosk_vendor_ba", "name": "Bà Ba Quầy Báo", "role": "Tiểu thương cửa ga", "location_id": "subway_kiosk", "influence": 0.40, "credibility": 0.45, "skepticism": 0.25, "sociality": 0.95, "deception_tendency": 0.18, "receptiveness": 0.85, "personality": "thích buôn chuyện, bàn tán tin giật gân", "private_motive": "Tụ tập đông người mua nước và bánh mì.", "direct_knowledge": "Thấy khói mờ bốc lên từ cửa thông gió ngầm.", "blind_spots": "Nói với mọi người là dưới ga bị nổ bom sập hầm rồi."},
        {"id": "s_livestreamer_long", "name": "Tiktoker Long", "role": "Livestreamer săn tin đường phố", "location_id": "street_entrance", "influence": 0.78, "credibility": 0.30, "skepticism": 0.20, "sociality": 0.95, "deception_tendency": 0.35, "receptiveness": 0.80, "personality": "kích động, câu view, đặt tiêu đề sốc", "private_motive": "Kiếm hàng trăm nghìn mắt xem livestream vụ 'thảm họa metro'.", "direct_knowledge": "Chỉ quay hình ảnh xe cảnh sát hú còi và người chạy nháo nhào lên từ cửa ga.", "blind_spots": "Bịa đặt rằng 'nghe tiếng nổ lớn và hàng trăm người nằm bất động dưới ga'."},
        {"id": "s_vent_engineer_an", "name": "Kỹ sư An", "role": "Vận hành Trạm Quạt Gió", "location_id": "ventilation_plant", "influence": 0.45, "credibility": 0.88, "skepticism": 0.80, "sociality": 0.30, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "kỹ thuật, chính xác", "private_motive": "Duy trì áp lực gió tươi 50m3/giây.", "direct_knowledge": "Cảm biến CO và độc chất trong luồng khí hoàn toàn ở mức an toàn 0 ppm.", "blind_spots": "Không lên trên mặt đất."},
        {"id": "s_passenger_tuan", "name": "Hành khách Tuấn", "role": "Kỹ sư xây dựng đi tàu E304", "location_id": "emergency_exit", "influence": 0.52, "credibility": 0.80, "skepticism": 0.70, "sociality": 0.60, "deception_tendency": 0.05, "receptiveness": 0.55, "personality": "tỉnh táo, biết quan sát", "private_motive": "Hỗ trợ phụ nữ và trẻ em leo cầu thang thoát hiểm an toàn.", "direct_knowledge": "Vỏ hầm bê tông nguyên vẹn, đèn khẩn cấp sáng rõ, không hề có sập hầm hay nứt vỡ.", "blind_spots": "Không biết tin đồn bên ngoài đã đi xa thế nào."},
        {"id": "s_security_guard_hung", "name": "Bảo vệ Hùng", "role": "Nhân viên an ninh sân ga", "location_id": "platform_1", "influence": 0.38, "credibility": 0.75, "skepticism": 0.65, "sociality": 0.50, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "chăm chỉ, trách nhiệm", "private_motive": "Đảm bảo không ai bị tụt lại trên đường ray.", "direct_knowledge": "Ray tiếp xúc thứ ba đã được ngắt điện hoàn toàn, hành khách đi bộ trật tự theo đèn chỉ dẫn.", "blind_spots": "Không biết chuyện gì ở các toa phía đuôi tàu."},
        {"id": "s_mother_huong", "name": "Chị Hương", "role": "Người nhà hành khách đi tàu", "location_id": "station_concourse", "influence": 0.35, "credibility": 0.45, "skepticism": 0.25, "sociality": 0.80, "deception_tendency": 0.05, "receptiveness": 0.90, "personality": "hoảng hốt, gào khóc gọi điện", "private_motive": "Tìm kiếm con gái đang trên chuyến tàu E304.", "direct_knowledge": "Điện thoại con gái mất sóng ngầm nên không liên lạc được.", "blind_spots": "Dễ dàng tin vào tin đồn tàu bị tấn công khí độc."},
        {"id": "s_journalist_quang", "name": "Nhà báo Quang", "role": "Phóng viên Giao thông Đô thị", "location_id": "street_entrance", "influence": 0.68, "credibility": 0.78, "skepticism": 0.70, "sociality": 0.70, "deception_tendency": 0.05, "receptiveness": 0.55, "personality": "chuyên nghiệp, kiểm chứng nguồn tin", "private_motive": "Đưa tin xác thực dập tắt tin đồn thất thiệt.", "direct_knowledge": "Đã phỏng vấn Chỉ huy Kiên và được xác nhận toàn bộ hành khách đã ra ngoài an toàn.", "blind_spots": "Chưa tiếp cận được báo cáo kỹ thuật nguyên nhân chập điện."},
        {"id": "s_rescue_worker_dung", "name": "Chiến sĩ Dũng", "role": "Lính Cứu nạn Hầm ngầm", "location_id": "train_tunnel_section", "influence": 0.40, "credibility": 0.85, "skepticism": 0.75, "sociality": 0.45, "deception_tendency": 0.02, "receptiveness": 0.40, "personality": "dũng cảm, tập trung cao độ", "private_motive": "Kiểm tra kỹ từng khoang tàu xem còn sót ai không.", "direct_knowledge": "Đã rà soát xong 6 toa tàu: 100% khoang trống, không có nạn nhân nào ngất xỉu hay tử vong.", "blind_spots": "Chỉ làm nhiệm vụ trong hầm."},
        {"id": "s_metro_director_binh", "name": "Giám đốc Bình", "role": "Phó Giám đốc Công ty Đường sắt Đô thị", "location_id": "station_control_room", "influence": 0.75, "credibility": 0.80, "skepticism": 0.75, "sociality": 0.45, "deception_tendency": 0.08, "receptiveness": 0.40, "personality": "thận trọng, cẩn trọng phát ngôn", "private_motive": "Báo cáo cấp trên và chuẩn bị họp báo thông tin chính thức.", "direct_knowledge": "Hệ thống tự động ATP hoạt động hoàn hảo đúng thiết kế an toàn khi gặp sự cố tín hiệu.", "blind_spots": "Chậm trễ trong việc phản hồi báo chí tại hiện trường."},
        {"id": "s_gate_attendant_lan", "name": "Nhân viên Lan", "role": "Soát vé Cổng Ga", "location_id": "ticket_gates", "influence": 0.32, "credibility": 0.70, "skepticism": 0.50, "sociality": 0.65, "deception_tendency": 0.05, "receptiveness": 0.60, "personality": "tận tâm, hòa nhã", "private_motive": "Giúp đỡ hành khách ra khỏi ga thông thoáng.", "direct_knowledge": "Thấy hành khách đi từ lối thoát hiểm lên đều đi lại bình thường, thở ổn định.", "blind_spots": "Không biết diễn biến trong buồng lái tàu."},
        {"id": "s_taxi_driver_thanh", "name": "Tài xế Thành", "role": "Lái xe Taxi trước cửa ga", "location_id": "street_entrance", "influence": 0.38, "credibility": 0.50, "skepticism": 0.30, "sociality": 0.85, "deception_tendency": 0.15, "receptiveness": 0.80, "personality": "thích chuyện giật gân, phát tán tin nhanh", "private_motive": "Tạo chủ đề nói chuyện với khách đi xe.", "direct_knowledge": "Thấy xe cứu hỏa bật đèn đỏ chớp nháy liên hồi.", "blind_spots": "Nghe người ta đồn có người chết ngạt là tin ngay."},
        {"id": "s_paramedic_viet", "name": "Điều dưỡng Việt", "role": "Cứu thương 115", "location_id": "field_triage_tent", "influence": 0.42, "credibility": 0.82, "skepticism": 0.75, "sociality": 0.50, "deception_tendency": 0.02, "receptiveness": 0.45, "personality": "chu đáo, tận tình", "private_motive": "Đo huyết áp và phát nước khoáng cho hành khách.", "direct_knowledge": "Không có ca nào phải chuyển viện cấp cứu hồi sức, tất cả đều tự đi lại được sau 10 phút nghỉ ngơi.", "blind_spots": "Không vào trong ga."}
    ]

    relationships = [
        {"from_agent_id": "s_driver_hai", "to_agent_id": "s_dispatcher_trang", "trust": 0.90, "reason": "Phối hợp điều hành lái tàu chuẩn mực"},
        {"from_agent_id": "s_dispatcher_trang", "to_agent_id": "s_tech_quoc", "trust": 0.88, "reason": "Xác nhận sự cố kỹ thuật tín hiệu"},
        {"from_agent_id": "s_fire_commander_kien", "to_agent_id": "s_dispatcher_trang", "trust": 0.92, "reason": "Chỉ huy cứu hộ liên ngành"},
        {"from_agent_id": "s_station_master_vinh", "to_agent_id": "s_police_son", "trust": 0.85, "reason": "Phối hợp an ninh nhà ga"},
        {"from_agent_id": "s_doctor_huyen", "to_agent_id": "s_fire_commander_kien", "trust": 0.88, "reason": "Bàn giao người cần sơ cứu"},
        {"from_agent_id": "s_passenger_nga", "to_agent_id": "s_doctor_huyen", "trust": 0.85, "reason": "Được bác sĩ chăm sóc trấn an"},
        {"from_agent_id": "s_kiosk_vendor_ba", "to_agent_id": "s_livestreamer_long", "trust": 0.75, "reason": "Kể chuyện giật gân cho streamer"},
        {"from_agent_id": "s_livestreamer_long", "to_agent_id": "s_taxi_driver_thanh", "trust": 0.70, "reason": "Bàn tán tin đồn giật gân"},
        {"from_agent_id": "s_mother_huong", "to_agent_id": "s_kiosk_vendor_ba", "trust": 0.75, "reason": "Hỏi thăm tình hình con gái"},
        {"from_agent_id": "s_journalist_quang", "to_agent_id": "s_fire_commander_kien", "trust": 0.90, "reason": "Lấy thông tin từ chỉ huy cứu hộ"},
        {"from_agent_id": "s_rescue_worker_dung", "to_agent_id": "s_driver_hai", "trust": 0.88, "reason": "Kiểm tra an toàn đoàn tàu"},
        {"from_agent_id": "s_vent_engineer_an", "to_agent_id": "s_dispatcher_trang", "trust": 0.85, "reason": "Báo cáo chỉ số không khí quạt gió"},
        {"from_agent_id": "s_passenger_tuan", "to_agent_id": "s_journalist_quang", "trust": 0.82, "reason": "Kể lại trải nghiệm thoát hiểm thực tế"},
        {"from_agent_id": "s_metro_director_binh", "to_agent_id": "s_station_master_vinh", "trust": 0.85, "reason": "Chỉ đạo điều hành nhà ga"}
    ]

    observations = [
        {"agent_id": "s_driver_hai", "claim_id": "claim_train_stopped_safely", "confidence": 0.99, "basis": "direct_observation"},
        {"agent_id": "s_driver_hai", "claim_id": "claim_all_evacuated_zero_deaths", "confidence": 0.90, "basis": "direct_observation"},
        {"agent_id": "s_tech_quoc", "claim_id": "claim_signal_cable_short", "confidence": 0.98, "basis": "hardware_inspection"},
        {"agent_id": "s_tech_quoc", "claim_id": "claim_no_structural_collapse", "confidence": 0.95, "basis": "direct_observation"},
        {"agent_id": "s_vent_engineer_an", "claim_id": "claim_ventilation_working", "confidence": 0.99, "basis": "telemetry_sensors"},
        {"agent_id": "s_fire_commander_kien", "claim_id": "claim_all_evacuated_zero_deaths", "confidence": 0.98, "basis": "operational_report"},
        {"agent_id": "s_doctor_huyen", "claim_id": "claim_all_evacuated_zero_deaths", "confidence": 0.95, "basis": "medical_triage_log"},
        {"agent_id": "s_passenger_tuan", "claim_id": "claim_no_structural_collapse", "confidence": 0.92, "basis": "direct_observation"}
    ]

    evidence = [
        {
            "id": "ev_occ_telemetry_log",
            "name": "Nhật ký vận hành tự động OCC Metro",
            "description": "Dữ liệu hộp đen điều độ: 'ATP triggered safe emergency stop on train E304 at 18:30:12 due to signal loss. Power grid isolated, tunnel airflow normal at 52 m3/s'.",
            "location_id": "station_control_room",
            "reliability": 0.99,
            "source_category": "technical_telemetry",
            "supports": ["claim_train_stopped_safely", "claim_ventilation_working"],
            "contradicts": ["claim_terrorist_gas_attack", "claim_hundreds_suffocated"],
            "unlock_condition": "Xem dữ liệu màn hình điều hành sau khi trao đổi với Điều độ viên Trang."
        },
        {
            "id": "ev_signal_box_debris",
            "name": "Mẫu cáp tín hiệu đường ray bị đoản mạch",
            "description": "Biên bản giám định kỹ thuật: Đoạn cáp tín hiệu hộp 14 bị chập cháy lớp vỏ PVC cách điện do rò rỉ nước ngầm tích tụ, không phát hiện dấu vết chất nổ hay tác động phá hoại.",
            "location_id": "track_signal_box",
            "reliability": 0.98,
            "source_category": "engineering_inspection",
            "supports": ["claim_signal_cable_short"],
            "contradicts": ["claim_terrorist_gas_attack"],
            "unlock_condition": "Khám nghiệm hộp tín hiệu sau khi phỏng vấn Kỹ sư Quốc."
        },
        {
            "id": "ev_fire_rescue_headcount",
            "name": "Biên bản hoàn thành cứu nạn cứu hộ",
            "description": "Biên bản chính thức của Cảnh sát PCCC & CNCH: Đã sơ tán toàn bộ 120 hành khách ra ngoài an toàn qua cửa thoát hiểm số 4; 0 người chết, 4 ca sơ cứu tâm lý tại chỗ.",
            "location_id": "street_entrance",
            "reliability": 0.99,
            "source_category": "official_headcount",
            "supports": ["claim_all_evacuated_zero_deaths", "claim_no_structural_collapse"],
            "contradicts": ["claim_hundreds_suffocated", "claim_tunnel_collapse_denial"],
            "unlock_condition": "Kiểm tra tại xe chỉ huy cứu hộ sau khi trao đổi với Trung tá Kiên."
        },
        {
            "id": "ev_tunnel_structural_scan",
            "name": "Bản quét siêu âm kết cấu hầm ngầm",
            "description": "Báo cáo kiểm tra trắc địa laser: Toàn bộ vỏ vòm bê tông đường hầm từ km 4+200 đến km 5+100 hoàn toàn nguyên vẹn, không có biến dạng hay sụt lún.",
            "location_id": "train_tunnel_section",
            "reliability": 0.97,
            "source_category": "structural_survey",
            "supports": ["claim_no_structural_collapse"],
            "contradicts": ["claim_tunnel_collapse_denial"],
            "unlock_condition": "Tìm thấy trong sổ tay kỹ thuật sau khi nói chuyện với Hành khách Tuấn hoặc Lái tàu Hải."
        }
    ]

    public_events = [
        {"game_second": 120, "event_type": "train_stopped", "headline": "Tàu Metro E304 dừng khẩn cấp giữa đường hầm", "description": "Tín hiệu đèn xanh chuyển đỏ tự động khiến hệ thống phanh kích hoạt.", "location_id": "train_tunnel_section"},
        {"game_second": 300, "event_type": "evacuation_started", "headline": "Lực lượng cứu nạn tiếp cận đưa hành khách đi bộ", "description": "Hành khách được hướng dẫn trật tự men theo lối đi an toàn trong hầm.", "location_id": "platform_1"},
        {"game_second": 480, "event_type": "terror_gas_rumor", "headline": "Tin đồn 'Khủng bố xả khí độc làm chết hàng trăm người'", "description": "Clip quay cảnh người dân hoảng sợ chạy lên từ cửa ga gây náo loạn mạng xã hội.", "location_id": "street_entrance"},
        {"game_second": 720, "event_type": "all_passengers_safe", "headline": "Toàn bộ 120 hành khách đã lên mặt đất an toàn", "description": "Lực lượng chức năng hoàn tất kiểm đếm hiện trường và lập trạm sơ cứu dã chiến.", "location_id": "emergency_exit"}
    ]

    truth_form = {
        "event_type": {
            "type": "enum",
            "label": "Bản chất sự cố đường sắt",
            "answer": "signal_failure",
            "options": [
                {"value": "signal_failure", "label": "Sự cố chập tín hiệu kích hoạt dừng tàu an toàn"},
                {"value": "terrorist_attack", "label": "Tấn công khủng bố bom khí độc"},
                {"value": "tunnel_collapse", "label": "Sập sụt lún kết cấu hầm ngầm"},
                {"value": "unknown", "label": "Chưa xác định"}
            ]
        },
        "cause_category": {
            "type": "enum",
            "label": "Nguyên nhân trực tiếp",
            "answer": "short_circuit_cable",
            "options": [
                {"value": "short_circuit_cable", "label": "Đoản mạch cáp tín hiệu hộp 14 do ẩm mốc"},
                {"value": "chemical_gas_device", "label": "Kích nổ thiết bị hơi ngạt độc hại"},
                {"value": "train_collision", "label": "Va chạm với phương tiện khác"},
                {"value": "unknown", "label": "Không rõ"}
            ]
        },
        "major_explosion": {
            "type": "boolean",
            "label": "Có vụ nổ phá hủy toa tàu hay sập hầm ngầm không?",
            "answer": False
        },
        "fatalities": {
            "type": "integer",
            "label": "Số người tử vong thực tế",
            "answer": 0,
            "min": 0,
            "max": 200
        },
        "management_issue": {
            "type": "enum",
            "label": "Trách nhiệm quản lý vận hành",
            "answer": "waterproofing_maintenance_delay",
            "options": [
                {"value": "waterproofing_maintenance_delay", "label": "Chậm xử lý chống thấm nước rò rỉ vào hộp cáp kỹ thuật"},
                {"value": "security_negligence", "label": "Bỏ lọt phần tử khủng bố mang chất cấm lên tàu"},
                {"value": "none", "label": "Không có sai phạm"}
            ]
        }
    }

    save_json(s_dir, "scenario.json", metadata)
    save_json(s_dir, "locations.json", locations)
    save_json(s_dir, "claims.json", claims)
    save_json(s_dir, "agents.json", agents)
    save_json(s_dir, "relationships.json", relationships)
    save_json(s_dir, "observations.json", observations)
    save_json(s_dir, "evidence.json", evidence)
    save_json(s_dir, "public-events.json", public_events)
    save_json(s_dir, "truth-form.json", truth_form)

# ==============================================================================
# SCENARIO 5: city-water-panic (Khủng Hoảng Nguồn Nước Thành Phố)
# ==============================================================================
def create_water_scenario():
    sc_id = "city-water-panic"
    s_dir = os.path.join(BASE_DIR, sc_id)

    metadata = {
        "id": sc_id,
        "title": "Khủng Hoảng Nguồn Nước Thành Phố",
        "version": 1,
        "description": "Lúc 06:15 sáng, van xả áp số 9 trên tuyến ống chính D800 của Trạm Cấp Nước Tây Thành Phố bị rách gioăng cao su, khiến dòng nước xoáy cuốn lớp phù sa oxit sắt tự nhiên sục lên làm nước máy khu dân cư đổi màu nâu đỏ trong 3 tiếng. Chỉ số độc tố hoàn toàn âm tính, 0 ca nhiễm độc. Tin đồn ác ý lan truyền rằng nguồn nước bị đầu độc bằng xyanua công nghiệp và đã có nhiều người tử vong.",
        "primary_false_narrative": {
            "adoption_threshold": 0.65,
            "defeat_ratio": 0.75,
            "claims": [
                {"claim_id": "claim_cyanide_poisoning", "weight": 0.35},
                {"claim_id": "claim_mass_poisoning_deaths", "weight": 0.40},
                {"claim_id": "claim_water_plant_coverup", "weight": 0.25}
            ]
        },
        "required_truth_fields": ["event_type", "cause_category", "major_explosion", "fatalities"],
        "bonus_truth_fields": ["water_toxicity", "management_issue"]
    }

    locations = [
        {"id": "water_plant_gate", "name": "Cổng Nhà Máy Nước", "description": "Người dân kéo đến mang theo xô chậu và mẫu nước đục yêu cầu giải thích.", "position": {"x": -24, "y": 0, "z": -16}, "connected_locations": ["pump_station", "water_quality_lab", "plant_admin"]},
        {"id": "pump_station", "name": "Trạm Bơm Áp Lực Cao (Tuyến D800)", "description": "Nơi đặt các máy bơm cao áp và van xả số 9 bị sự cố áp lực.", "position": {"x": -32, "y": 0, "z": -24}, "connected_locations": ["water_plant_gate", "pipeline_junction"]},
        {"id": "pipeline_junction", "name": "Hố Van Tuyến Ống Chính", "description": "Vị trí van xả áp lực số 9 bị rách gioăng kỹ thuật.", "position": {"x": -28, "y": 0, "z": -8}, "connected_locations": ["pump_station"]},
        {"id": "water_quality_lab", "name": "Phòng Thử Nghiệm Chất Lượng Nước", "description": "Phòng kiểm nghiệm hóa sinh tiến hành phân tích quang phổ chỉ tiêu nước.", "position": {"x": -20, "y": 0, "z": -26}, "connected_locations": ["water_plant_gate", "filtration_tanks"]},
        {"id": "filtration_tanks", "name": "Hệ Thống Bể Lọc Cát & Khử Trùng", "description": "Các bể lọc cát thạch anh và clo hóa nước sạch.", "position": {"x": -14, "y": 0, "z": -18}, "connected_locations": ["water_quality_lab"]},
        {"id": "plant_admin", "name": "Văn Phòng Ban Lãnh Đạo", "description": "Ban giám đốc họp khẩn để cô lập phân đoạn ống đục.", "position": {"x": 2, "y": 0, "z": 6}, "connected_locations": ["water_plant_gate", "supermarket"]},
        {"id": "supermarket", "name": "Siêu Thị Trung Tâm", "description": "Người dân tranh nhau vét sạch toàn bộ nước khoáng đóng chai.", "position": {"x": 26, "y": 0, "z": -20}, "connected_locations": ["plant_admin", "residential_quarter"]},
        {"id": "residential_quarter", "name": "Khu Đô Thị Sông Xanh", "description": "Khu chung cư nơi cư dân phát hiện nước vòi chảy ra màu nâu đỏ.", "position": {"x": -6, "y": 0, "z": -6}, "connected_locations": ["supermarket", "local_clinic", "riverside_teahouse"]},
        {"id": "local_clinic", "name": "Trạm Y Tế Phường", "description": "Nơi người dân kéo đến đòi xét nghiệm máu vì sợ trúng độc.", "position": {"x": 10, "y": 0, "z": -4}, "connected_locations": ["residential_quarter"]},
        {"id": "riverside_teahouse", "name": "Quán Trà Bờ Sông", "description": "Điểm bàn luận rôm rả của cư dân về màu nước máy buổi sáng.", "position": {"x": 18, "y": 0, "z": 12}, "connected_locations": ["residential_quarter", "environment_office"]},
        {"id": "environment_office", "name": "Sở Tài Nguyên & Môi Trường", "description": "Thanh tra môi trường lấy mẫu nước độc lập tại nguồn cấp.", "position": {"x": -22, "y": 0, "z": 14}, "connected_locations": ["riverside_teahouse"]},
        {"id": "police_station", "name": "Công An Khu Vực", "description": "Chốt công an giữ trật tự trước cảnh gom nước đóng chai tại siêu thị.", "position": {"x": -10, "y": 0, "z": 18}, "connected_locations": ["supermarket"]}
    ]

    claims = [
        {"id": "claim_water_discoloration_occurred", "text": "Nước máy sinh hoạt tại khu Tây bị đục và đổi màu nâu đỏ trong 3 giờ", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_mineral_sediment_origin", "text": "Hiện tượng do lớp bùn oxit sắt khoáng tự nhiên bị sục lên khi van áp lực gặp sự cố", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_toxicity_negative", "text": "Mẫu nước kiểm nghiệm âm tính 100% với độc tố xyanua và kim loại nặng nguy hiểm", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_zero_fatalities", "text": "Không có bất kỳ trường hợp nào ngộ độc hay tử vong vì nước máy", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_water_flushed_clear", "text": "Hệ thống đã xả đáy xong và nước máy đã trong trở lại đạt tiêu chuẩn", "classification": "true", "role": "fact", "weight": 0.0},
        {"id": "claim_cyanide_poisoning", "text": "Nguồn nước sạch thành phố đã bị đầu độc bằng hóa chất xyanua công nghiệp", "classification": "false", "role": "primary_false", "weight": 0.35},
        {"id": "claim_mass_poisoning_deaths", "text": "Nhiều gia đình uống nước buổi sáng đã bị co giật tử vong tại nhà", "classification": "false", "role": "primary_false", "weight": 0.40},
        {"id": "claim_water_plant_coverup", "text": "Công ty nước sạch biết có độc từ tuần trước nhưng cố tình giấu nhẹm", "classification": "false", "role": "primary_false", "weight": 0.25},
        {"id": "claim_toxic_river_waste", "text": "Có tàu chở rác thải lén đổ hóa chất độc hại vào bể chứa đầu nguồn", "classification": "false", "role": "secondary_false", "weight": 0.0},
        {"id": "claim_citywide_water_cutoff_month", "text": "Toàn thành phố sẽ bị cắt nước sinh hoạt trong 1 tháng tới", "classification": "false", "role": "secondary_false", "weight": 0.0}
    ]

    agents = [
        {"id": "w_chief_eng_thanh", "name": "Kỹ sư trưởng Thành", "role": "Giám đốc Kỹ thuật Nhà máy Nước", "location_id": "pump_station", "influence": 0.70, "credibility": 0.88, "skepticism": 0.80, "sociality": 0.35, "deception_tendency": 0.05, "receptiveness": 0.40, "personality": "khoa học, điềm tĩnh, hiểu tường tận mạng lưới", "private_motive": "Khắc phục van xả và mở xả đáy làm sạch đường ống.", "direct_knowledge": "Van xả áp số 9 rách gioăng khiến áp lực tụt từ 6 bar xuống 2 bar, đảo chiều dòng chảy cuốn theo cáu cặn oxit sắt bám thành ống.", "blind_spots": "Không ngờ tin đồn lại biến thành xyanua đầu độc."},
        {"id": "w_chemist_quynh", "name": "Dược sĩ Quỳnh", "role": "Trưởng phòng Thử nghiệm Nước", "location_id": "water_quality_lab", "influence": 0.65, "credibility": 0.95, "skepticism": 0.88, "sociality": 0.30, "deception_tendency": 0.02, "receptiveness": 0.30, "personality": "chính xác tuyệt đối, trung thực", "private_motive": "Công bố kết quả thử nghiệm quang phổ mẫu nước minh bạch.", "direct_knowledge": "Hàm lượng sắt đạt 0.6 mg/L (cao hơn bình thường gây đục), nhưng chỉ số xyanua, asen, thủy ngân hoàn toàn bằng 0.", "blind_spots": "Bận làm xét nghiệm nên chưa thông tin ra ngoài sảnh."},
        {"id": "w_director_phat", "name": "Giám đốc Phát", "role": "Tổng Giám đốc Công ty Cấp nước", "location_id": "plant_admin", "influence": 0.80, "credibility": 0.75, "skepticism": 0.75, "sociality": 0.40, "deception_tendency": 0.10, "receptiveness": 0.35, "personality": "chính khách, thận trọng, lo ngại truyền thông", "private_motive": "Bảo vệ hợp đồng nhượng quyền cấp nước của công ty.", "direct_knowledge": "Đã cho ngừng bơm tuyến D800 và cấp bù nước sạch từ nhà máy dự phòng phía Đông.", "blind_spots": "Chần chừ không tổ chức họp báo sớm."},
        {"id": "w_operator_duc", "name": "Công nhân Đức", "role": "Thợ vận hành Hố van", "location_id": "pipeline_junction", "influence": 0.40, "credibility": 0.80, "skepticism": 0.65, "sociality": 0.50, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "chăm chỉ, cần cù", "private_motive": "Thay gioăng cao su mới để khóa van an toàn.", "direct_knowledge": "Đã tháo van số 9 ra, gioăng cao su chịu lực bị nứt rách do lão hóa tự nhiên.", "blind_spots": "Không biết tin tức trên mạng xã hội."},
        {"id": "w_doctor_mai", "name": "Bác sĩ Mai", "role": "Trưởng Trạm Y Tế Phường", "location_id": "local_clinic", "influence": 0.60, "credibility": 0.90, "skepticism": 0.85, "sociality": 0.55, "deception_tendency": 0.02, "receptiveness": 0.40, "personality": "chu đáo, tận tình, có uy tín với dân cư", "private_motive": "Khám và trấn an người dân đến trạm.", "direct_knowledge": "Từ sáng tới giờ tiếp nhận 15 người đến khám do lo sợ, nhưng không ai có triệu chứng ngộ độc xyanua (không khó thở, không tím tái, mạch bình thường).", "blind_spots": "Không có thiết bị kiểm tra mẫu nước tại chỗ."},
        {"id": "w_citizen_bac_nam", "name": "Bác Nam Tổ Trưởng", "role": "Tổ trưởng Dân phố Sông Xanh", "location_id": "residential_quarter", "influence": 0.55, "credibility": 0.65, "skepticism": 0.35, "sociality": 0.85, "deception_tendency": 0.05, "receptiveness": 0.85, "personality": "lo cho cộng đồng nhưng dễ bị kích động", "private_motive": "Bảo vệ tính mạng bà con trong tổ dân phố.", "direct_knowledge": "Hứng xô nước thấy màu đỏ quạch như gạch cua, để lắng 15 phút thì thấy cặn nâu dưới đáy xô.", "blind_spots": "Tin vào bài đăng trên Zalo bảo có người ngộ độc chết ở tòa nhà kế bên."},
        {"id": "w_supermarket_mgr", "name": "Quản lý Hùng", "role": "Quản lý Siêu thị Trung tâm", "location_id": "supermarket", "influence": 0.50, "credibility": 0.68, "skepticism": 0.50, "sociality": 0.70, "deception_tendency": 0.12, "receptiveness": 0.60, "personality": "thực tế, kinh doanh nhanh nhạy", "private_motive": "Bán sạch hàng tồn kho nước đóng chai nhưng sợ vỡ trận chen lấn.", "direct_knowledge": "Toàn bộ 3.000 thùng nước tinh khiết đóng chai bán sạch trong vòng 40 phút.", "blind_spots": "Không nắm được khi nào nước máy phục hồi."},
        {"id": "w_env_inspector_tram", "name": "Thanh tra Trâm", "role": "Thanh tra Sở Tài Nguyên Môi Trường", "location_id": "environment_office", "influence": 0.70, "credibility": 0.92, "skepticism": 0.88, "sociality": 0.40, "deception_tendency": 0.02, "receptiveness": 0.30, "personality": "độc lập, công minh, khắt khe", "private_motive": "Lấy mẫu nước kiểm tra chéo độc lập để xử phạt nếu có sai phạm.", "direct_knowledge": "Mẫu nước lấy tại đầu vòi khu dân cư có nồng độ oxy hòa tan và pH bình thường, cặn nâu chủ yếu là Fe2O3 vô hại.", "blind_spots": "Đang hoàn tất văn bản pháp lý chính thức."},
        {"id": "w_blogger_trieu", "name": "Blogger Triệu", "role": "Reviewer Đời sống & Phốt", "location_id": "riverside_teahouse", "influence": 0.75, "credibility": 0.35, "skepticism": 0.20, "sociality": 0.95, "deception_tendency": 0.30, "receptiveness": 0.80, "personality": "thích giật tít, phóng đại, tạo phốt", "private_motive": "Tăng tương tác cho fanpage cá nhân.", "direct_knowledge": "Chỉ lấy bức ảnh xô nước đỏ của cư dân rồi viết bài giật tít: 'Báo động: Xyanua tràn vào vòi nước, đã có ca tử vong!'", "blind_spots": "Hoàn toàn bịa đặt thông tin về các ca tử vong."},
        {"id": "w_guard_tuan", "name": "Bảo vệ Tuấn", "role": "Bảo vệ Cổng Nhà máy Nước", "location_id": "water_plant_gate", "influence": 0.35, "credibility": 0.60, "skepticism": 0.45, "sociality": 0.70, "deception_tendency": 0.05, "receptiveness": 0.65, "personality": "vất vả, cố gắng giải thích", "private_motive": "Ngăn không cho người dân quá khích xông vào khu xử lý nước.", "direct_knowledge": "Các kỹ sư đang thay van và xả sục rửa đường ống từ sáng sớm.", "blind_spots": "Không rành về thuật ngữ hóa học."},
        {"id": "w_police_kien", "name": "Trung úy Kiên", "role": "Cảnh sát Trật tự Phường", "location_id": "police_station", "influence": 0.55, "credibility": 0.85, "skepticism": 0.80, "sociality": 0.40, "deception_tendency": 0.02, "receptiveness": 0.35, "personality": "nghiêm khắc, bảo đảm trật tự", "private_motive": "Dẹp cảnh tranh giành găm hàng nước uống tại siêu thị.", "direct_knowledge": "Công an phường đã kiểm tra toàn bộ địa bàn: không có vụ án mạng hay ca ngộ độc chết người nào.", "blind_spots": "Không trả lời được câu hỏi kỹ thuật về nước."},
        {"id": "w_teahouse_owner_hoa", "name": "Bà Hoa Trà Đá", "role": "Chủ quán nước ven sông", "location_id": "riverside_teahouse", "influence": 0.45, "credibility": 0.48, "skepticism": 0.25, "sociality": 0.90, "deception_tendency": 0.15, "receptiveness": 0.85, "personality": "nhiều chuyện, thêm mắm dặm muối", "private_motive": "Kể chuyện ly kỳ để giữ chân khách uống nước.", "direct_knowledge": "Sáng nay pha trà bằng nước máy thấy trà đổi sang màu đen tím ngắt (phản ứng tự nhiên giữa tanin trong trà và ion sắt).", "blind_spots": "Tưởng màu đen là do thuốc độc hóa học."},
        {"id": "w_mother_lan", "name": "Chị Lan Nội Trợ", "role": "Cư dân Khu Sông Xanh", "location_id": "residential_quarter", "influence": 0.38, "credibility": 0.55, "skepticism": 0.30, "sociality": 0.80, "deception_tendency": 0.05, "receptiveness": 0.85, "personality": "lo sợ cho con cái", "private_motive": "Không dám nấu cơm bằng nước máy.", "direct_knowledge": "Đã phải dắt con sang nhà người quen ở quận khác tắm nhờ.", "blind_spots": "Hốt hoảng vì đọc phải tin đồn thất thiệt của Blogger Triệu."},
        {"id": "w_water_tech_ha", "name": "Kỹ thuật viên Hà", "role": "Vận hành Bể Lọc Cát", "location_id": "filtration_tanks", "influence": 0.42, "credibility": 0.85, "skepticism": 0.75, "sociality": 0.40, "deception_tendency": 0.02, "receptiveness": 0.40, "personality": "cần mẫn, chu đáo", "private_motive": "Đảm bảo các bể lọc cát thạch anh hoạt động tốt.", "direct_knowledge": "Nước sau bể lọc tại nhà máy trong vắt 100%, đục chỉ xảy ra ở đường ống ngoài do sục cặn ống cũ.", "blind_spots": "Không tham gia xử lý ngoài mạng lưới truyền tải."},
        {"id": "w_journalist_phuong", "name": "Phóng viên Phương", "role": "Nhà báo Môi trường & Sức khỏe", "location_id": "environment_office", "influence": 0.68, "credibility": 0.80, "skepticism": 0.72, "sociality": 0.65, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "khách quan, đi sâu vào thực tế", "private_motive": "Viết bài phân tích khoa học để bác bỏ tin đồn đầu độc.", "direct_knowledge": "Đã phỏng vấn Dược sĩ Quỳnh và Thanh tra Trâm, nắm chắc kết quả xét nghiệm âm tính với chất độc.", "blind_spots": "Cần chờ biên bản có mộc đỏ chính thức."},
        {"id": "w_plumber_vinh", "name": "Thợ sửa ống nước Vinh", "role": "Thợ sửa ống nước độc lập", "location_id": "residential_quarter", "influence": 0.45, "credibility": 0.72, "skepticism": 0.60, "sociality": 0.75, "deception_tendency": 0.05, "receptiveness": 0.60, "personality": "thực tế, quen tay nghề", "private_motive": "Đi thông tắc và xả cặn cho các hộ gia đình.", "direct_knowledge": "Mở vòi xả hết 2 xô nước đục thì nước bắt đầu trong vắt trở lại bình thường.", "blind_spots": "Chỉ làm việc ở quy mô hộ gia đình."},
        {"id": "w_store_clerk_thuy", "name": "Thu ngân Thủy", "role": "Thu ngân Siêu thị", "location_id": "supermarket", "influence": 0.32, "credibility": 0.60, "skepticism": 0.40, "sociality": 0.70, "deception_tendency": 0.05, "receptiveness": 0.65, "personality": "mệt mỏi vì khách chen lấn", "private_motive": "Tính tiền nhanh cho khách.", "direct_knowledge": "Nhiều người mua cả chục lốc nước khoáng vì nghe tin đồn mất nước cả tháng.", "blind_spots": "Không biết chuyện gì ở nhà máy nước."},
        {"id": "w_student_manh", "name": "Sinh viên Mạnh", "role": "Sinh viên Khoa Hóa", "location_id": "riverside_teahouse", "influence": 0.42, "credibility": 0.75, "skepticism": 0.70, "sociality": 0.60, "deception_tendency": 0.05, "receptiveness": 0.55, "personality": "thích giải thích khoa học", "private_motive": "Giải thích hiện tượng đổi màu cho các bác lớn tuổi.", "direct_knowledge": "Màu đỏ gạch là rỉ sắt Fe2O3 bị xới tung lên khi thay đổi áp lực dòng chảy, hoàn toàn không phải xyanua.", "blind_spots": "Ít người chịu lắng nghe phân tích hóa học."},
        {"id": "w_driver_bao", "name": "Tài xế Bảo", "role": "Lái xe bồn cấp nước sạch", "location_id": "water_plant_gate", "influence": 0.38, "credibility": 0.75, "skepticism": 0.60, "sociality": 0.55, "deception_tendency": 0.05, "receptiveness": 0.50, "personality": "nhiệt tình, chu đáo", "private_motive": "Chở nước sạch khẩn cấp đến phục vụ miễn phí cho khu chung cư.", "direct_knowledge": "Xe bồn lấy nước trực tiếp từ trạm bơm phía Đông, nước đạt chuẩn 100%.", "blind_spots": "Không vào trong hầm van."},
        {"id": "w_resident_duc", "name": "Bác Đức Cựu Chiến Binh", "role": "Cư dân lâu năm", "location_id": "residential_quarter", "influence": 0.50, "credibility": 0.80, "skepticism": 0.75, "sociality": 0.55, "deception_tendency": 0.02, "receptiveness": 0.45, "personality": "chững chạc, bình tĩnh", "private_motive": "Khuyên bà con không nghe lời kẻ xấu kích động.", "direct_knowledge": "Nhắc lại rằng 5 năm trước nhà máy sửa đường ống cũng từng xảy ra hiện tượng nước đỏ tương tự.", "blind_spots": "Không có điện thoại thông minh để đọc tin tức mới."}
    ]

    relationships = [
        {"from_agent_id": "w_chief_eng_thanh", "to_agent_id": "w_chemist_quynh", "trust": 0.90, "reason": "Cộng sự kiểm soát kỹ thuật và chất lượng nước"},
        {"from_agent_id": "w_chief_eng_thanh", "to_agent_id": "w_operator_duc", "trust": 0.88, "reason": "Chỉ đạo trực tiếp thợ sửa hố van"},
        {"from_agent_id": "w_director_phat", "to_agent_id": "w_chief_eng_thanh", "trust": 0.85, "reason": "Tin cậy kỹ sư trưởng"},
        {"from_agent_id": "w_chemist_quynh", "to_agent_id": "w_env_inspector_tram", "trust": 0.92, "reason": "Đồng nghiệp khoa học kiểm định mẫu nước"},
        {"from_agent_id": "w_doctor_mai", "to_agent_id": "w_chemist_quynh", "trust": 0.85, "reason": "Xác minh chỉ số độc chất phục vụ y tế"},
        {"from_agent_id": "w_citizen_bac_nam", "to_agent_id": "w_doctor_mai", "trust": 0.88, "reason": "Tin tưởng bác sĩ trạm y tế"},
        {"from_agent_id": "w_teahouse_owner_hoa", "to_agent_id": "w_blogger_trieu", "trust": 0.75, "reason": "Kể chuyện nước đổi màu cho blogger"},
        {"from_agent_id": "w_blogger_trieu", "to_agent_id": "w_citizen_bac_nam", "trust": 0.60, "reason": "Lan truyền bài viết trên mạng xã hội"},
        {"from_agent_id": "w_police_kien", "to_agent_id": "w_doctor_mai", "trust": 0.85, "reason": "Xác minh tình hình an ninh sức khỏe"},
        {"from_agent_id": "w_journalist_phuong", "to_agent_id": "w_env_inspector_tram", "trust": 0.88, "reason": "Lấy kết quả kiểm định chính thức"},
        {"from_agent_id": "w_student_manh", "to_agent_id": "w_teahouse_owner_hoa", "trust": 0.70, "reason": "Khách quen uống nước giải thích khoa học"},
        {"from_agent_id": "w_plumber_vinh", "to_agent_id": "w_citizen_bac_nam", "trust": 0.80, "reason": "Thợ sửa nước tin cậy của xóm"},
        {"from_agent_id": "w_driver_bao", "to_agent_id": "w_chief_eng_thanh", "trust": 0.85, "reason": "Nhận lệnh điều động xe bồn"},
        {"from_agent_id": "w_resident_duc", "to_agent_id": "w_citizen_bac_nam", "trust": 0.85, "reason": "Hàng xóm cựu chiến binh"}
    ]

    observations = [
        {"agent_id": "w_chief_eng_thanh", "claim_id": "claim_water_discoloration_occurred", "confidence": 0.99, "basis": "direct_observation"},
        {"agent_id": "w_chief_eng_thanh", "claim_id": "claim_mineral_sediment_origin", "confidence": 0.95, "basis": "pressure_sensor_logs"},
        {"agent_id": "w_chemist_quynh", "claim_id": "claim_toxicity_negative", "confidence": 0.99, "basis": "spectrometry_assay"},
        {"agent_id": "w_chemist_quynh", "claim_id": "claim_zero_fatalities", "confidence": 0.90, "basis": "toxicology_report"},
        {"agent_id": "w_doctor_mai", "claim_id": "claim_zero_fatalities", "confidence": 0.98, "basis": "clinical_examination_log"},
        {"agent_id": "w_env_inspector_tram", "claim_id": "claim_toxicity_negative", "confidence": 0.95, "basis": "independent_sample_test"},
        {"agent_id": "w_operator_duc", "claim_id": "claim_mineral_sediment_origin", "confidence": 0.92, "basis": "valve_inspection"},
        {"agent_id": "w_plumber_vinh", "claim_id": "claim_water_flushed_clear", "confidence": 0.90, "basis": "tap_water_flushing"}
    ]

    evidence = [
        {
            "id": "ev_water_spectrometry_report",
            "name": "Phiếu phân tích quang phổ mẫu nước (Lab Analysis)",
            "description": "Kết quả phân tích mẫu nước: Sắt Fe tổng = 0.58 mg/L (vượt nhẹ do rỉ sét ống cũ), Xyanua CN- < 0.001 mg/L (âm tính tuyệt đối), Asen < 0.001 mg/L, Clo dư 0.4 mg/L. Nước hoàn toàn không chứa độc tố nguy hiểm.",
            "location_id": "water_quality_lab",
            "reliability": 0.99,
            "source_category": "chemical_analysis",
            "supports": ["claim_toxicity_negative", "claim_mineral_sediment_origin"],
            "contradicts": ["claim_cyanide_poisoning", "claim_mass_poisoning_deaths"],
            "unlock_condition": "Kiểm tra phòng xét nghiệm sau khi phỏng vấn Dược sĩ Quỳnh."
        },
        {
            "id": "ev_valve_pressure_log",
            "name": "Nhật ký áp lực trạm bơm & hiện trường van 9",
            "description": "Bản ghi áp lực đường ống: Van xả áp số 9 rách gioăng lúc 06:15 gây sụt áp đột ngột tạo dòng xoáy sục cáu cặn oxit sắt bám đáy ống. Đã cô lập van và thay gioăng hoàn tất lúc 08:30.",
            "location_id": "pump_station",
            "reliability": 0.98,
            "source_category": "engineering_record",
            "supports": ["claim_mineral_sediment_origin", "claim_water_discoloration_occurred"],
            "contradicts": ["claim_cyanide_poisoning"],
            "unlock_condition": "Khám nghiệm trạm bơm sau khi nói chuyện với Kỹ sư trưởng Thành hoặc Công nhân Đức."
        },
        {
            "id": "ev_clinic_toxicology_clearance",
            "name": "Sổ theo dõi khám bệnh Trạm Y Tế",
            "description": "Báo cáo y tế phường: Tiếp nhận 15 người dân, 100% bệnh nhân sinh hiệu bình thường, không có triệu chứng ngộ độc xyanua; 0 ca cấp cứu chuyển tuyến, 0 ca tử vong.",
            "location_id": "local_clinic",
            "reliability": 0.99,
            "source_category": "medical_clearance",
            "supports": ["claim_zero_fatalities"],
            "contradicts": ["claim_mass_poisoning_deaths"],
            "unlock_condition": "Xem hồ sơ tại trạm y tế sau khi gặp Bác sĩ Mai."
        },
        {
            "id": "ev_env_independent_certificate",
            "name": "Chứng thư kiểm định độc lập Sở TN&MT",
            "description": "Biên bản kiểm tra độc lập tại 5 điểm cấp nước khu dân cư xác nhận các chỉ số lý hóa đã trở lại ngưỡng an toàn chuẩn QCVN 01-1:2018/BYT sau khi hệ thống xả đáy.",
            "location_id": "environment_office",
            "reliability": 0.98,
            "source_category": "environmental_audit",
            "supports": ["claim_water_flushed_clear", "claim_toxicity_negative"],
            "contradicts": ["claim_water_plant_coverup", "claim_citywide_water_cutoff_month"],
            "unlock_condition": "Lấy tại văn phòng Sở TN&MT sau khi trao đổi với Thanh tra Trâm."
        }
    ]

    public_events = [
        {"game_second": 120, "event_type": "water_discolored", "headline": "Cư dân phát hiện nước máy chảy ra màu nâu đỏ", "description": "Nhiều hộ gia đình mở vòi nước sinh hoạt buổi sáng thấy nước đục ngầu.", "location_id": "residential_quarter"},
        {"game_second": 300, "event_type": "water_panic_buy", "headline": "Cơn sốt vét sạch nước đóng chai tại siêu thị", "description": "Hàng trăm người xếp hàng dài tranh nhau mua nước khoáng đóng bình.", "location_id": "supermarket"},
        {"game_second": 480, "event_type": "cyanide_rumor", "headline": "Tin đồn ác ý: 'Nguồn nước nhiễm xyanua gây chết người'", "description": "Bài đăng thất thiệt lan truyền gây tâm lý hoang mang tột độ cho người dân đô thị.", "location_id": "riverside_teahouse"},
        {"game_second": 720, "event_type": "flush_completed", "headline": "Nhà máy nước hoàn tất xả đáy sục rửa đường ống", "description": "Dòng nước tại các vòi xả đã trong trẻo trở lại và đạt tiêu chuẩn an toàn.", "location_id": "water_plant_gate"}
    ]

    truth_form = {
        "event_type": {
            "type": "enum",
            "label": "Bản chất sự cố nguồn nước",
            "answer": "pipe_rupture",
            "options": [
                {"value": "pipe_rupture", "label": "Sự cố van áp lực làm sục cặn oxit sắt gây đục tạm thời"},
                {"value": "chemical_poisoning", "label": "Nguồn nước bị đầu độc bằng xyanua công nghiệp"},
                {"value": "algal_bloom", "label": "Bùng phát tảo độc ô nhiễm nguồn nước"},
                {"value": "unknown", "label": "Chưa xác định"}
            ]
        },
        "cause_category": {
            "type": "enum",
            "label": "Nguyên nhân chính",
            "answer": "mineral_sediment_surge",
            "options": [
                {"value": "mineral_sediment_surge", "label": "Dòng xoáy sục lớp bùn oxit sắt tự nhiên trong thành ống cũ"},
                {"value": "deliberate_sabotage", "label": "Kẻ xấu lén đổ hóa chất độc hại vào bể chứa"},
                {"value": "industrial_effluent", "label": "Nhà máy xả trộm nước thải chưa qua xử lý"},
                {"value": "unknown", "label": "Không rõ"}
            ]
        },
        "major_explosion": {
            "type": "boolean",
            "label": "Có thảm họa nhiễm độc diện rộng / nguồn nước bị nhiễm độc chết người không?",
            "answer": False
        },
        "fatalities": {
            "type": "integer",
            "label": "Số người tử vong thực tế",
            "answer": 0,
            "min": 0,
            "max": 50
        },
        "management_issue": {
            "type": "enum",
            "label": "Trách nhiệm quản lý",
            "answer": "aging_gasket_replacement_delay",
            "options": [
                {"value": "aging_gasket_replacement_delay", "label": "Chậm thay thế gioăng cao su chịu lực định kỳ của van xả"},
                {"value": "falsified_safety_records", "label": "Làm giả phiếu kiểm nghiệm chất lượng nước sạch"},
                {"value": "none", "label": "Không có sai phạm"}
            ]
        }
    }

    save_json(s_dir, "scenario.json", metadata)
    save_json(s_dir, "locations.json", locations)
    save_json(s_dir, "claims.json", claims)
    save_json(s_dir, "agents.json", agents)
    save_json(s_dir, "relationships.json", relationships)
    save_json(s_dir, "observations.json", observations)
    save_json(s_dir, "evidence.json", evidence)
    save_json(s_dir, "public-events.json", public_events)
    save_json(s_dir, "truth-form.json", truth_form)

if __name__ == "__main__":
    print("Generating 4 scenarios...")
    create_hospital_scenario()
    create_bank_scenario()
    create_subway_scenario()
    create_water_scenario()
    print("All 4 scenarios generated successfully!")
