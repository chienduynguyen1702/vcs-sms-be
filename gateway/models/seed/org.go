package seed

import (
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/models"
)

var Org = models.Organization{
	Name:              "Viettel Cyber Security",
	AliasName:         "VCS",
	Description:       "Viettel Cyber Security (Công ty An ninh mạng Viettel) là đơn vị trực thuộc Tập đoàn Công nghiệp - Viễn thông Quân đội với nhiệm vụ chính là thực hiện toàn trình từ nghiên cứu chuyên sâu, phát triển các giải pháp về ATTT trên môi trường viễn thông, CNTT đến cung cấp các sản phẩm/dịch vụ ATTT do chính Công ty xây dựng tới khách hàng lớn trong và ngoài nước. Viettel Cyber Security tự hào có đội ngũ nhân sự giỏi, những chuyên gia hàng đầu trong lĩnh vực Công nghệ thông tin nói chung và ATTT nói riêng. Với môi trường làm việc trẻ trung, năng động và hiện đại, Viettel Cyber Security luôn mong muốn có thể thu hút được nhiều nhân tài gia nhập vào đại gia đình Viettel.",
	WebsiteURL:        "https://viettelcybersecurity.com/",
	Address:           "Keangnam Hanoi Landmark Tower, Mễ Trì, Từ Liêm, Hà Nội",
	EstablishmentDate: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
}
