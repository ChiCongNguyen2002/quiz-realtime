# Real-Time Quiz System

## Bài toán là gì?

Thử tưởng tượng bạn đang tổ chức một cuộc thi quiz online. Có hàng trăm người cùng chơi một lúc.

Bạn muốn:
- Mỗi khi một người trả lời đúng, TẤT CẢ người chơi khác thấy ngay là ai đang dẫn đầu
- Không phải chờ ai bấm refresh gì cả
- Phải thật nhanh, gần như ngay lập tức

Đó chính là bài toán mà hệ thống này giải quyết.

---

## Yêu cầu cụ thể là gì?

1. **User Participation** - Người chơi có thể join vào quiz session bằng session ID. Nhiều người cùng join vào một session được.

2. **Real-Time Score Updates** - Khi người chơi submit đáp án, điểm số được cập nhật ngay lập tức. Hệ thống chấm điểm phải chính xác và nhất quán.

3. **Real-Time Leaderboard** - Leaderboard hiển thị thứ hạng của tất cả người tham gia. Cập nhật ngay khi có ai đó thay đổi điểm.

---

## Có những cách nào để giải quyết bài toán này?

### Cách 1: HTTP Polling (Gọi API liên tục)

Đây là cách đơn giản nhất mà nhiều người nghĩ đến đầu tiên.

**Hoạt động thế nào:**
- Trình duyệt cứ 1-2 giây gọi một lần lên server: "Ê, có ai mới submit chưa?"
- Server trả về danh sách điểm hiện tại
- Trình duyệt hiển thị lên màn hình

**Ưu điểm:**
- Dễ viết, ai cũng làm được
- Không cần công nghệ gì phức tạp

**Nhược điểm:**
- Có độ trễ 1-2 giây (thực ra cũng khá là chậm)
- Tốn băng thông vô ích (gọi liên tục dù không có gì mới)
- Server chịu tải lớn khi nhiều người
- Màn hình nhảy liên tục, nhìn khó chịu

---

### Cách 2: WebSocket

Đây là cách mà hệ thống này đang dùng.

**Hoạt động thế nào:**
- Trình duyệt mở kết nối WebSocket với server, giữ kết nối đó suốt
- Khi muốn submit đáp án, gọi API bình thường
- Server xử lý xong, đẩy kết quả qua WebSocket cho TẤT CẢ người trong phòng

**Ưu điểm:**
- Cực nhanh (gần như ngay lập tức)
- Kết nối 2 chiều
- Tốn ít tài nguyên
- Hỗ trợ hàng ngàn người cùng lúc

**Nhược điểm:**
- Phức tạp hơn HTTP thông thường
- Cần quản lý kết nối (nhưng đã có thư viện lo)

---

## Giải pháp hiện tại hoạt động thế nào?

### Tổng quan

Hệ thống gồm 4 phần chính:
- **API Server (Go/Gin)**: Nhận request, xử lý nghiệp vụ
- **WebSocket Hub**: Quản lý kết nối, gửi tin nhắn
- **PostgreSQL**: Lưu câu hỏi, kết quả (database)
- **Redis**: Lưu leaderboard (siêu nhanh)

### Flow chi tiết từng bước

**Bước 1: Người chơi vào phòng**

Người chơi mở trình duyệt, kết nối WebSocket đến server, gửi kèm session_id để server biết người này thuộc phòng nào.

Server nhận kết nối, lưu vào "phòng" tương ứng.

**Bước 2: Người chơi trả lời câu hỏi**

Người chơi gửi đáp án qua API:
- Server nhận được đáp án
- Lấy câu hỏi từ database
- So sánh đáp án của user với đáp án đúng
- Tính tổng điểm (đúng mỗi câu +1 điểm)

**Bước 3: Lưu kết quả**

- Lưu điểm vào PostgreSQL (để lưu trữ lâu dài)
- Cập nhật leaderboard trong Redis (để truy vấn nhanh)

**Bước 4: Gửi kết quả cho TẤT CẢ người trong phòng**

Đây là bước quan trọng nhất!

- Server tìm tất cả những người đang ở trong cùng session
- Gửi tin nhắn WebSocket cho họ
- Tin nhắn chứa: ai vừa submit, được bao nhiêu điểm, bảng xếp hạng mới

**Bước 5: Người chơi nhận được tin**

Trình duyệt nhận được tin nhắn WebSocket, hiển thị leaderboard mới lên màn hình.

Người chơi thấy ngay: "À, thằng A vừa được 8 điểm, nó đang dẫn đầu!"

---

## Giải pháp này giải quyết được những gì?

| Yêu cầu | Có giải quyết không? |
|----------|----------------|
| User join session |  OK |
| Nhiều user cùng session |  OK |
| Real-time score |  OK |
| Real-time leaderboard |  OK |
| Scale được |  OK |

---

## Tổng kết

Bài toán này về cơ bản là bài toán **truyền tin real-time**.

Giải pháp hiện tại:
- Dùng WebSocket để gửi tin nhắn tức thì
- Redis để xử lý leaderboard nhanh nhất có thể
- Clean architecture để code dễ đọc, dễ test, dễ bảo trì
