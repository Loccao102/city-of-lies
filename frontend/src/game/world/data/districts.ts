import type { LocationInfo } from "@/types/game";

export interface DistrictLocation extends LocationInfo {
  color: number;
}

export const DISTRICT_LOCATIONS: DistrictLocation[] = [
  { id: "factory_gate", name: "Cổng Nhà Máy", description: "Lối vào chính của khu Riverside.", x: -24, z: -16, color: 0x4a5568, category: "factory" },
  { id: "storage_building", name: "Nhà Kho DB-4", description: "Tâm điểm của sự cố và hiện trường cần kiểm tra.", x: -32, z: -24, color: 0xe53e3e, category: "factory" },
  { id: "security_room", name: "Phòng An Ninh", description: "Camera, nhật ký ra vào và nhân viên bảo vệ.", x: -28, z: -8, color: 0x2b6cb0, category: "factory" },
  { id: "factory_yard", name: "Sân Nhà Máy", description: "Không gian tập trung công nhân và lực lượng phản ứng.", x: -20, z: -26, color: 0x718096, category: "factory" },
  { id: "factory_office", name: "Văn Phòng Quản Lý", description: "Nơi làm việc của ban quản lý nhà máy.", x: -14, z: -18, color: 0xd69e2e, category: "factory" },
  { id: "hospital", name: "Bệnh Viện Đa Khoa", description: "Nguồn xác minh tình trạng thương vong và y tế.", x: 26, z: -20, color: 0x319795, category: "official" },
  { id: "fire_station", name: "Trạm Cứu Hỏa", description: "Nguồn thông tin chính thức về đám cháy và phản ứng khẩn cấp.", x: -6, z: -6, color: 0xc53030, category: "official" },
  { id: "news_office", name: "Tòa Soạn & Blogger", description: "Điểm khuếch đại thông tin công khai và tin đồn.", x: 18, z: 12, color: 0x805ad5, category: "public" },
  { id: "public_square", name: "Quảng Trường", description: "Không gian công cộng có mật độ trao đổi thông tin cao.", x: 2, z: 6, color: 0xdd6b20, category: "public" },
  { id: "market", name: "Chợ Dân Sinh", description: "Mạng lưới truyền miệng dân cư hoạt động mạnh.", x: -10, z: 18, color: 0x38a169, category: "public" },
  { id: "cafe", name: "Quán Cà Phê Vỉa Hè", description: "Điểm gặp gỡ không chính thức của nhân chứng và người dân.", x: 10, z: -4, color: 0xd69e2e, category: "public" },
  { id: "residential_street", name: "Khu Dân Cư", description: "Khu vực chịu tác động trực tiếp từ tin đồn lan truyền.", x: -22, z: 14, color: 0x4a5568, category: "residential" },
];

export const DISTRICT_LOCATION_BY_ID = new Map(
  DISTRICT_LOCATIONS.map((location) => [location.id, location] as const)
);
