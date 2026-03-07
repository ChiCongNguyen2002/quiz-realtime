# Real-Time Quiz System

Đây là một hệ thống quiz realtime đơn giản. Người dùng có thể:

- Tham gia một **session quiz**
- Trả lời câu hỏi
- Thấy **điểm và bảng xếp hạng cập nhật ngay**

Hệ thống sử dụng:

- **Go** cho backend
- **Postgres** để lưu dữ liệu chính
- **Redis** để tăng tốc leaderboard
- **WebSocket** để cập nhật realtime

---

# 1. Bài toán

Hệ thống cần giải quyết ba vấn đề chính.

## User Participation
Người dùng phải có thể **tham gia một session quiz bằng mã session**.  
Nhiều người có thể cùng tham gia một session và thi với nhau.

## Real-Time Score Updates
Khi người chơi gửi đáp án:

- Hệ thống phải **tính điểm ngay**
- Điểm phải **chính xác và không bị trùng**

## Real-Time Leaderboard
Session phải có **bảng xếp hạng của tất cả người chơi**.  
Leaderboard cần **cập nhật ngay khi điểm thay đổi**.

---

# 2. Cách hệ thống hoạt động

## Tạo session
Host tạo một session cho một bộ câu hỏi.  
Hệ thống tạo ra một **session ID** để người chơi tham gia.

---

## Người chơi tham gia session
Người chơi nhập **session ID** để vào phòng.

Hệ thống lưu lại thông tin:

- session
- user

Một session có thể có **nhiều người chơi cùng lúc**.

---

## Người chơi nộp đáp án
Khi người chơi submit bài:

Server sẽ:

1. Lấy danh sách câu hỏi của quiz
2. So sánh đáp án
3. Tính điểm
4. Lưu điểm vào **Postgres**
5. Cập nhật điểm vào **Redis leaderboard**

Mỗi user trong một session chỉ có **một bản ghi điểm**.

---

## Leaderboard realtime
Sau khi điểm được cập nhật:

- Leaderboard của session được cập nhật
- Server gửi leaderboard mới tới tất cả người chơi qua **WebSocket**

Người chơi sẽ thấy bảng xếp hạng **thay đổi ngay trên màn hình**.

---

## Khi realtime không hoạt động
Nếu WebSocket gặp lỗi, client vẫn có thể lấy leaderboard bằng API.

Dữ liệu luôn được lưu trong **Postgres**, nên không bị mất.

---

# 3. Các thành phần chính

## Session
Một session đại diện cho **một lần chơi quiz**.

Một quiz có thể có **nhiều session khác nhau**.

---

## Participant
Danh sách người chơi trong session.

Mỗi participant gắn với:

- session
- user

---

## Score
Điểm của user trong session.

Score được lưu trong **Postgres** để đảm bảo dữ liệu an toàn.

---

## Leaderboard
Leaderboard hiển thị thứ hạng người chơi trong session.

Redis được dùng để:

- cập nhật điểm nhanh
- lấy danh sách top người chơi nhanh

---

# 4. Kiến trúc

Hệ thống được chia thành các phần chính:

- **Domain** – chứa entity và logic nghiệp vụ
- **Application** – xử lý các luồng chính như join session, submit answer
- **Infrastructure** – làm việc với Postgres, Redis và WebSocket
- **Delivery** – cung cấp HTTP API

Cách chia này giúp code **dễ đọc, dễ test và dễ mở rộng**.

---

# 5. Chạy project

## Chạy local

Yêu cầu:

- Go
- Postgres
- Redis

Chạy server:


- go run ./cmd/server


Server mặc định chạy ở port **8080**.

---

## Chạy bằng Docker

Build image:

- make docker-build

Run container:

- make docker-run

Sau khi run là có thể tiến hành test và apply đúng như yêu cầu bài toán trên thông qua các tool như Postman