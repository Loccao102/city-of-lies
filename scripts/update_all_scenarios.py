import json
import os

BASE_DIR = os.path.join(os.path.dirname(__file__), "..", "data", "scenarios")

def calibrate_scenario(sc_id, initial_obs_to_add, events):
    s_dir = os.path.join(BASE_DIR, sc_id)
    obs_path = os.path.join(s_dir, "observations.json")
    with open(obs_path, "r", encoding="utf-8") as f:
        obs = json.load(f)
    
    # filter out existing speculative seeds if any to avoid duplicates
    existing_claims = {o["claim_id"] for o in initial_obs_to_add}
    obs = [o for o in obs if not (o.get("basis") in ("speculation", "hearsay") and o["claim_id"] in existing_claims)]
    obs.extend(initial_obs_to_add)
    
    with open(obs_path, "w", encoding="utf-8") as f:
        json.dump(obs, f, ensure_ascii=False, indent=2)
    
    pe_path = os.path.join(s_dir, "public-events.json")
    with open(pe_path, "w", encoding="utf-8") as f:
        json.dump(events, f, ensure_ascii=False, indent=2)
    print(f"Successfully calibrated {sc_id}")

# 1. midtown-bank-run
bank_obs = [
    {"agent_id": "b_broker_trung", "claim_id": "claim_bank_insolvency", "confidence": 0.40, "basis": "speculation"},
    {"agent_id": "b_vendor_dung", "claim_id": "claim_bank_insolvency", "confidence": 0.35, "basis": "hearsay"},
    {"agent_id": "b_citizen_ba", "claim_id": "claim_bank_insolvency", "confidence": 0.30, "basis": "hearsay"}
]

bank_events = [
    {"game_second": 120, "event_type": "atm_glitch", "headline": "Các cây ATM Midtown đồng loạt báo lỗi kết nối", "description": "Hàng trăm khách hàng không thể rút tiền mặt hoặc kiểm tra số dư qua ứng dụng di động.", "location_id": "atm_gallery", "belief_seeds": []},
    {"game_second": 300, "event_type": "crowd_gathering", "headline": "Đám đông kéo đến đại sảnh ngân hàng", "description": "Người dân lo lắng tập trung yêu cầu rút tiền tiết kiệm trước hạn.", "location_id": "bank_hall", "belief_seeds": []},
    {"game_second": 420, "event_type": "rumor_seed", "headline": "Tài xế Dũng bàn tán về việc ngân hàng hết tiền", "description": "Dũng nói với bà Ba rằng thấy nhân viên tắt đèn cây ATM.", "location_id": "atm_gallery", "belief_seeds": [
        {"agent_id": "b_vendor_dung", "claim_id": "claim_bank_insolvency", "confidence": 0.45, "reasoning": "Thấy bảo vệ xua tay báo cây ATM ngừng giao dịch."},
        {"agent_id": "b_citizen_ba", "claim_id": "claim_bank_insolvency", "confidence": 0.40, "reasoning": "Nghe tài xế bảo máy chủ ngân hàng bị sập."}
    ]},
    {"game_second": 480, "event_type": "rumor_transmission", "headline": "Khách VIP lo ngại mất tiền gửi", "description": "Tin đồn hệ thống máy chủ bị tê liệt lan truyền.", "location_id": "vip_lounge", "belief_seeds": [
        {"agent_id": "b_citizen_ba", "claim_id": "claim_bank_insolvency", "confidence": 0.50, "reasoning": "Càng lúc càng đông người kéo đến rút tiền."},
        {"agent_id": "b_investor_lam", "claim_id": "claim_bank_insolvency", "confidence": 0.45, "reasoning": "App VIP Banking không đăng nhập được."}
    ]},
    {"game_second": 720, "event_type": "blog_post", "headline": "Môi giới Trung đăng tin cảnh báo thanh khoản", "description": "Bài viết: Ngân hàng Midtown tê liệt mạng, nghi mất thanh khoản.", "location_id": "broker_cafe", "belief_seeds": [
        {"agent_id": "b_broker_trung", "claim_id": "claim_bank_insolvency", "confidence": 0.60, "reasoning": "Cổ phiếu ngân hàng bị bán tháo sàn."}
    ]},
    {"game_second": 840, "event_type": "headline_update", "headline": "Trung đổi tiêu đề: 'Khủng hoảng thanh khoản Midtown'", "description": "Bài đăng cảnh báo người dân khẩn trương rút tiền.", "location_id": "broker_cafe", "belief_seeds": [
        {"agent_id": "b_broker_trung", "claim_id": "claim_bank_insolvency", "confidence": 0.72, "reasoning": "Thấy tin tức chia sẻ chóng mặt trên các room chứng khoán."},
        {"agent_id": "b_reporter_tuan", "claim_id": "claim_bank_insolvency", "confidence": 0.68, "reasoning": "Nhiều nguồn tin đồn ngân hàng sắp bị kiểm soát đặc biệt."}
    ]},
    {"game_second": 900, "event_type": "social_viral", "headline": "Làn sóng rút tiền lan rộng trên mạng xã hội", "description": "Các hội nhóm chia sẻ thông tin khiến đám đông kéo tới xếp hàng dài.", "location_id": "financial_square", "belief_seeds": [
        {"agent_id": "b_investor_lam", "claim_id": "claim_bank_insolvency", "confidence": 0.80, "reasoning": "Nhận được cuộc gọi cảnh báo từ đối tác kinh doanh."},
        {"agent_id": "b_reporter_tuan", "claim_id": "claim_bank_insolvency", "confidence": 0.75, "reasoning": "Chưa thấy thông cáo báo chí giải thích thỏa đáng."}
    ]},
    {"game_second": 1080, "event_type": "official_statement", "headline": "Lãnh đạo ngân hàng chậm giải thích số liệu thanh khoản", "description": "Ban điều hành chưa công bố báo cáo kiểm toán khiến tin đồn phong tỏa bùng nổ.", "location_id": "boardroom", "belief_seeds": [
        {"agent_id": "b_broker_trung", "claim_id": "claim_account_freeze_order", "confidence": 0.65, "reasoning": "Tin đồn lệnh phong tỏa tài khoản sẽ ban hành lúc 17:00."},
        {"agent_id": "b_citizen_ba", "claim_id": "claim_account_freeze_order", "confidence": 0.65, "reasoning": "Người xếp hàng bảo sau chiều nay sẽ không rút được nữa."},
        {"agent_id": "b_customer_lan", "claim_id": "claim_account_freeze_order", "confidence": 0.60, "reasoning": "Lệnh chuyển khoản công ty vẫn bị treo."},
        {"agent_id": "b_vendor_dung", "claim_id": "claim_account_freeze_order", "confidence": 0.60, "reasoning": "Thấy bảo vệ đóng cửa sắt chặn khách."},
        {"agent_id": "b_investor_lam", "claim_id": "claim_account_freeze_order", "confidence": 0.65, "reasoning": "Luật sư cảnh báo rủi ro phong tỏa tư pháp."},
        {"agent_id": "b_reporter_tuan", "claim_id": "claim_account_freeze_order", "confidence": 0.70, "reasoning": "Báo chí đặt câu hỏi về tính pháp lý giao dịch."},
        {"agent_id": "b_security_guard", "claim_id": "claim_account_freeze_order", "confidence": 0.50, "reasoning": "Nhận lệnh giữ khoảng cách với người dân."},
        {"agent_id": "b_teller_mai", "claim_id": "claim_account_freeze_order", "confidence": 0.48, "reasoning": "Màn hình hệ thống bị khóa chưa mở lại."}
    ]},
    {"game_second": 1200, "event_type": "livestream_claim", "headline": "Môi giới Trung livestream khẳng định Tổng giám đốc đã tháo chạy", "description": "Trung nói trên sóng trực tiếp rằng xe chở vàng đã bí mật rời kho.", "location_id": "broker_cafe", "belief_seeds": [
        {"agent_id": "b_broker_trung", "claim_id": "claim_ceo_fled_with_gold", "confidence": 0.75, "reasoning": "Khẳng định có tin nội bộ lãnh đạo đã bay sang nước ngoài."},
        {"agent_id": "b_vendor_dung", "claim_id": "claim_ceo_fled_with_gold", "confidence": 0.68, "reasoning": "Khách bàn tán sếp ngân hàng tẩu tán tài sản."},
        {"agent_id": "b_citizen_ba", "claim_id": "claim_ceo_fled_with_gold", "confidence": 0.55, "reasoning": "Khóc lóc vì tưởng mất hết tiền tiết kiệm dưỡng già."}
    ]},
    {"game_second": 1320, "event_type": "headline_fatalities", "headline": "Các trang tin lá cải đăng bài: 'Sếp ngân hàng ôm vàng bỏ trốn'", "description": "Tin đồn thất thiệt đạt đỉnh điểm chia sẻ trên các nền tảng mạng.", "location_id": "broker_cafe", "belief_seeds": [
        {"agent_id": "b_broker_trung", "claim_id": "claim_ceo_fled_with_gold", "confidence": 0.82, "reasoning": "Cập nhật bài viết khẳng định vàng dự trữ đã bị chuyển đi."},
        {"agent_id": "b_reporter_tuan", "claim_id": "claim_ceo_fled_with_gold", "confidence": 0.78, "reasoning": "Áp lực tin nóng đẩy nghi vấn lên trang nhất."}
    ]}
]

# 2. subway-line3-standstill
subway_obs = [
    {"agent_id": "s_kiosk_vendor_ba", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.35, "basis": "hearsay"},
    {"agent_id": "s_livestreamer_long", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.40, "basis": "speculation"},
    {"agent_id": "s_taxi_driver_thanh", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.30, "basis": "hearsay"}
]

subway_events = [
    {"game_second": 120, "event_type": "train_stopped", "headline": "Tàu Metro E304 dừng khẩn cấp giữa đường hầm", "description": "Tín hiệu đèn xanh chuyển đỏ tự động khiến hệ thống phanh kích hoạt.", "location_id": "train_tunnel_section", "belief_seeds": []},
    {"game_second": 300, "event_type": "evacuation_started", "headline": "Lực lượng cứu nạn tiếp cận đưa hành khách đi bộ", "description": "Hành khách được hướng dẫn trật tự men theo lối đi an toàn trong hầm.", "location_id": "platform_1", "belief_seeds": []},
    {"game_second": 420, "event_type": "rumor_seed", "headline": "Bà Ba quầy báo bàn tán về khói bốc lên", "description": "Bà Ba nói với tài xế Thành rằng ngửi thấy mùi khét độc hại từ cửa thông gió.", "location_id": "subway_kiosk", "belief_seeds": [
        {"agent_id": "s_kiosk_vendor_ba", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.45, "reasoning": "Thấy khói mờ bốc lên từ hầm thông gió kèm mùi khét lẹt."},
        {"agent_id": "s_taxi_driver_thanh", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.40, "reasoning": "Nghe tiếng còi hụ xe cứu thương dồn dập."}
    ]},
    {"game_second": 480, "event_type": "rumor_transmission", "headline": "Thân nhân hành khách hoang mang", "description": "Chị Hương gào khóc tìm con khi thấy xe cứu hỏa đổ về ga.", "location_id": "station_concourse", "belief_seeds": [
        {"agent_id": "s_taxi_driver_thanh", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.50, "reasoning": "Hành khách chạy lên từ cửa ga mặt mày tái mét."},
        {"agent_id": "s_mother_huong", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.45, "reasoning": "Mất liên lạc điện thoại với con gái trong hầm."}
    ]},
    {"game_second": 720, "event_type": "blog_post", "headline": "Tiktoker Long livestream nghi vấn khủng bố", "description": "Video quay cảnh xe hú còi giật tít: 'Báo động đỏ: Nghi nổ bom ga ngầm'.", "location_id": "street_entrance", "belief_seeds": [
        {"agent_id": "s_livestreamer_long", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.60, "reasoning": "Quay cảnh lực lượng phản ứng nhanh mang mặt nạ phòng độc."}
    ]},
    {"game_second": 840, "event_type": "headline_update", "headline": "Long đổi tiêu đề: 'Khủng bố xả hơi ngạt trên tàu Metro?'", "description": "Bài đăng cảnh báo hơi độc nhận hàng chục nghìn lượt xem trực tiếp.", "location_id": "street_entrance", "belief_seeds": [
        {"agent_id": "s_livestreamer_long", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.72, "reasoning": "Có người bình luận nói thấy hành khách ngất xỉu nôn mửa."},
        {"agent_id": "s_mother_huong", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.68, "reasoning": "Xem livestream hốt hoảng khóc ngất."}
    ]},
    {"game_second": 900, "event_type": "social_viral", "headline": "Mạng xã hội hoảng loạn chia sẻ tin tấn công khí độc", "description": "Các diễn đàn giao thông lan truyền thông tin thất thiệt.", "location_id": "street_entrance", "belief_seeds": [
        {"agent_id": "s_mother_huong", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.80, "reasoning": "Thấy đám đông nhốn nháo lan truyền bài viết."},
        {"agent_id": "s_journalist_quang", "claim_id": "claim_terrorist_gas_attack", "confidence": 0.75, "reasoning": "Tin tức lan quá nhanh chưa kịp kiểm chứng hiện trường."}
    ]},
    {"game_second": 1080, "event_type": "official_statement", "headline": "Giám đốc đường sắt né tránh xác nhận tình trạng đường hầm", "description": "Phát biểu quanh co khiến dư luận nghi ngờ sập hầm và bưng bít thông tin.", "location_id": "station_control_room", "belief_seeds": [
        {"agent_id": "s_livestreamer_long", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.65, "reasoning": "Lãnh đạo metro không công bố hình ảnh bên trong hầm."},
        {"agent_id": "s_kiosk_vendor_ba", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.65, "reasoning": "Nghe đồn trần hầm bị nứt sập kẹt tàu."},
        {"agent_id": "s_taxi_driver_thanh", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.60, "reasoning": "Thấy xe cứu hộ mang thiết bị khoan cắt vào hầm."},
        {"agent_id": "s_mother_huong", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.60, "reasoning": "Lo sợ con gái bị kẹt dưới đống đổ nát."},
        {"agent_id": "s_passenger_nga", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.60, "reasoning": "Lúc chạy thấy bụi bay mù mịt trong hầm tối."},
        {"agent_id": "s_journalist_quang", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.70, "reasoning": "Nhà ga đóng cửa kín mít không cho báo chí vào."},
        {"agent_id": "s_station_master_vinh", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.50, "reasoning": "Áp lực điều phối khiến căng thẳng cao độ."},
        {"agent_id": "s_gate_attendant_lan", "claim_id": "claim_tunnel_collapse_denial", "confidence": 0.48, "reasoning": "Khách thoát hiểm lên kể hầm rung lắc mạnh."}
    ]},
    {"game_second": 1200, "event_type": "livestream_claim", "headline": "Tiktoker Long livestream khẳng định hàng trăm người chết ngạt", "description": "Long phát sóng nói có người thân gọi báo nạn nhân ngất xỉu hàng loạt.", "location_id": "street_entrance", "belief_seeds": [
        {"agent_id": "s_livestreamer_long", "claim_id": "claim_hundreds_suffocated", "confidence": 0.75, "reasoning": "Khẳng định trên livestream có hàng trăm người chết ngạt trong toa kín."},
        {"agent_id": "s_kiosk_vendor_ba", "claim_id": "claim_hundreds_suffocated", "confidence": 0.68, "reasoning": "Nghe người ta kháo nhau xe cứu thương chở xác đi."},
        {"agent_id": "s_mother_huong", "claim_id": "claim_hundreds_suffocated", "confidence": 0.55, "reasoning": "Nghe tin có người chết ngạt liền ngất lịm."}
    ]},
    {"game_second": 1320, "event_type": "headline_fatalities", "headline": "Báo lá cải đưa tin thảm họa thương vong đường hầm", "description": "Thông tin chưa kiểm chứng được lan truyền với tốc độ chóng mặt.", "location_id": "street_entrance", "belief_seeds": [
        {"agent_id": "s_livestreamer_long", "claim_id": "claim_hundreds_suffocated", "confidence": 0.82, "reasoning": "Cập nhật tiêu đề xác nhận thảm kịch chết ngạt hàng loạt."},
        {"agent_id": "s_journalist_quang", "claim_id": "claim_hundreds_suffocated", "confidence": 0.78, "reasoning": "Bài viết câu view của đối thủ cạnh tranh leo top thịnh hành."}
    ]}
]

# 3. city-water-panic
water_obs = [
    {"agent_id": "w_blogger_trieu", "claim_id": "claim_cyanide_poisoning", "confidence": 0.40, "basis": "speculation"},
    {"agent_id": "w_teahouse_owner_hoa", "claim_id": "claim_cyanide_poisoning", "confidence": 0.35, "basis": "hearsay"},
    {"agent_id": "w_mother_lan", "claim_id": "claim_cyanide_poisoning", "confidence": 0.30, "basis": "hearsay"}
]

water_events = [
    {"game_second": 120, "event_type": "water_discolored", "headline": "Cư dân phát hiện nước máy chảy ra màu nâu đỏ", "description": "Nhiều hộ gia đình mở vòi nước sinh hoạt buổi sáng thấy nước đục ngầu.", "location_id": "residential_quarter", "belief_seeds": []},
    {"game_second": 300, "event_type": "water_panic_buy", "headline": "Cơn sốt vét sạch nước đóng chai tại siêu thị", "description": "Hàng trăm người xếp hàng dài tranh nhau mua nước khoáng đóng bình.", "location_id": "supermarket", "belief_seeds": []},
    {"game_second": 420, "event_type": "rumor_seed", "headline": "Quán trà bà Hoa râm ran tin nước đổi màu đen tím", "description": "Bà Hoa nói với bác Nam rằng nước máy pha trà chuyển màu đen kỳ dị nghi có độc.", "location_id": "riverside_teahouse", "belief_seeds": [
        {"agent_id": "w_teahouse_owner_hoa", "claim_id": "claim_cyanide_poisoning", "confidence": 0.45, "reasoning": "Nước máy pha trà chuyển sang màu tím đen bất thường."},
        {"agent_id": "w_citizen_bac_nam", "claim_id": "claim_cyanide_poisoning", "confidence": 0.40, "reasoning": "Hứng nước thấy cặn đỏ quạch như gạch cua bốc mùi tanh."}
    ]},
    {"game_second": 480, "event_type": "rumor_transmission", "headline": "Cư dân bàn tán nghi vấn ô nhiễm hóa chất", "description": "Người dân dặn nhau tuyệt đối không dùng nước máy nấu ăn.", "location_id": "residential_quarter", "belief_seeds": [
        {"agent_id": "w_citizen_bac_nam", "claim_id": "claim_cyanide_poisoning", "confidence": 0.50, "reasoning": "Cả khu chung cư nhốn nháo khóa van nước tổng."},
        {"agent_id": "w_mother_lan", "claim_id": "claim_cyanide_poisoning", "confidence": 0.45, "reasoning": "Sợ con uống phải nước độc bị co giật."}
    ]},
    {"game_second": 720, "event_type": "blog_post", "headline": "Blogger Triệu đăng ảnh xô nước đỏ cảnh báo độc chất", "description": "Bài viết: Nghi vấn đầu độc nguồn nước sinh hoạt phía Tây.", "location_id": "riverside_teahouse", "belief_seeds": [
        {"agent_id": "w_blogger_trieu", "claim_id": "claim_cyanide_poisoning", "confidence": 0.60, "reasoning": "Nhận được hàng trăm tin nhắn ảnh nước đỏ từ người dân."}
    ]},
    {"game_second": 840, "event_type": "headline_update", "headline": "Triệu đổi tiêu đề: 'Thảm họa xyanua tràn vào vòi nước'", "description": "Bài đăng cảnh báo chất kịch độc xyanua đạt hàng chục nghìn lượt chia sẻ hoảng loạn.", "location_id": "riverside_teahouse", "belief_seeds": [
        {"agent_id": "w_blogger_trieu", "claim_id": "claim_cyanide_poisoning", "confidence": 0.72, "reasoning": "Có tin đồn công ty xả hóa chất xyanua chưa qua xử lý."},
        {"agent_id": "w_mother_lan", "claim_id": "claim_cyanide_poisoning", "confidence": 0.68, "reasoning": "Đọc tin thấy trùng khớp với triệu chứng nước đổi màu đen."}
    ]},
    {"game_second": 900, "event_type": "social_viral", "headline": "Làn sóng gom nước sạch gây náo loạn siêu thị", "description": "Cư dân chia sẻ tin đồn khiến các cửa hàng cạn kiệt nước tinh khiết.", "location_id": "supermarket", "belief_seeds": [
        {"agent_id": "w_mother_lan", "claim_id": "claim_cyanide_poisoning", "confidence": 0.80, "reasoning": "Tranh cướp bình nước đóng chai ở siêu thị với cư dân khác."},
        {"agent_id": "w_journalist_phuong", "claim_id": "claim_cyanide_poisoning", "confidence": 0.75, "reasoning": "Thông tin đầu độc nguồn nước lan tràn trên tất cả các mạng."}
    ]},
    {"game_second": 1080, "event_type": "official_statement", "headline": "Ban giám đốc nhà máy nước từ chối công bố nguyên nhân", "description": "Lãnh đạo chậm trễ họp báo giải thích hiện tượng nước đục khiến dư luận phẫn nộ.", "location_id": "plant_admin", "belief_seeds": [
        {"agent_id": "w_blogger_trieu", "claim_id": "claim_water_plant_coverup", "confidence": 0.65, "reasoning": "Nhà máy đóng cổng không cho phóng viên tiếp cận bể lọc."},
        {"agent_id": "w_teahouse_owner_hoa", "claim_id": "claim_water_plant_coverup", "confidence": 0.65, "reasoning": "Bà con kéo đến cổng nhà máy đòi câu trả lời thỏa đáng."},
        {"agent_id": "w_citizen_bac_nam", "claim_id": "claim_water_plant_coverup", "confidence": 0.60, "reasoning": "Tổ dân phố kiến nghị công ty giải trình nhưng bị từ chối."},
        {"agent_id": "w_mother_lan", "claim_id": "claim_water_plant_coverup", "confidence": 0.60, "reasoning": "Nghi ngờ công ty bưng bít sự cố hóa chất."},
        {"agent_id": "w_supermarket_mgr", "claim_id": "claim_water_plant_coverup", "confidence": 0.55, "reasoning": "Khách hàng kháo nhau nước máy sẽ bị ngắt dài hạn."},
        {"agent_id": "w_journalist_phuong", "claim_id": "claim_water_plant_coverup", "confidence": 0.70, "reasoning": "Sở Tài nguyên Môi trường chưa nhận được báo cáo từ nhà máy."},
        {"agent_id": "w_guard_tuan", "claim_id": "claim_water_plant_coverup", "confidence": 0.50, "reasoning": "Lãnh đạo chỉ đạo khóa cổng bảo vệ."},
        {"agent_id": "w_store_clerk_thuy", "claim_id": "claim_water_plant_coverup", "confidence": 0.48, "reasoning": "Nghe khách kể công ty giấu tin xyanua."},
        {"agent_id": "w_plumber_vinh", "claim_id": "claim_water_plant_coverup", "confidence": 0.45, "reasoning": "Thợ sửa ống không được cấp phụ tùng thay thế."}
    ]},
    {"game_second": 1200, "event_type": "livestream_claim", "headline": "Blogger Triệu livestream báo động nhiều người co giật tử vong", "description": "Triệu khẳng định đã có ca cấp cứu ngộ độc xyanua tử vong tại nhà.", "location_id": "riverside_teahouse", "belief_seeds": [
        {"agent_id": "w_blogger_trieu", "claim_id": "claim_mass_poisoning_deaths", "confidence": 0.75, "reasoning": "Khẳng định có gia đình uống nước buổi sáng đã bị co giật tử vong."},
        {"agent_id": "w_teahouse_owner_hoa", "claim_id": "claim_mass_poisoning_deaths", "confidence": 0.68, "reasoning": "Khách uống trà bàn tán có người đưa đi cấp cứu không qua khỏi."},
        {"agent_id": "w_mother_lan", "claim_id": "claim_mass_poisoning_deaths", "confidence": 0.55, "reasoning": "Hốt hoảng ôm con khóc vì sợ trúng độc."}
    ]},
    {"game_second": 1320, "event_type": "headline_fatalities", "headline": "Các trang mạng giật tít: 'Khủng hoảng đầu độc nước sạch cướp đi sinh mạng'", "description": "Tin giả đạt đỉnh điểm tác động tiêu cực tới toàn bộ người dân thành phố.", "location_id": "riverside_teahouse", "belief_seeds": [
        {"agent_id": "w_blogger_trieu", "claim_id": "claim_mass_poisoning_deaths", "confidence": 0.82, "reasoning": "Cập nhật bài viết xác nhận nhiều ca tử vong do độc tố."},
        {"agent_id": "w_journalist_phuong", "claim_id": "claim_mass_poisoning_deaths", "confidence": 0.78, "reasoning": "Áp lực thông tin khiến chuẩn bị đăng bài phản ánh phản ứng dư luận."}
    ]}
]

if __name__ == "__main__":
    calibrate_scenario("midtown-bank-run", bank_obs, bank_events)
    calibrate_scenario("subway-line3-standstill", subway_obs, subway_events)
    calibrate_scenario("city-water-panic", water_obs, water_events)
    print("All scenarios updated and calibrated!")
