from locust import HttpUser, task, between
from random import randrange

test_data = [
    'Kapan Diklat diadakan oleh OSIS Seksi Bidang 4, dan siapa yang bisa mengikuti sebagai peserta?',
    'Berapa frekuensi pelaksanaan "Pembiasaan Senam Pagi" setiap bulannya?',
    'Apa yang menjadi fokus kegiatan "Bulan Bahasa" di Kampus SMAN SITURAJA untuk Seksi Bidang 8?',
    'Mengapa pembinaan kesiswaan perlu mendukung kreativitas, keterampilan, dan kewirausahaan?',
    'Apa tujuan dari program Memperingati Maulid Nabi Muhammad SAW?',
    'Mengapa kondisi yang tertib, keamanan terjaga, dan situasi kondusif dianggap esensial dalam pendidikan?',
    'Bagaimana Wakil Kepala Bidang Humas berkoordinasi dengan instansi pemerintah dan dinas terkait?',
    'Apa kegiatan tahunan yang diselenggarakan oleh Wakil Kepala Bidang Humas untuk seluruh warga SMAN Situraja?',
    'Bagaimana sekolah mengukur aspek Kreatif dalam mencapai visinya?',
    'Siapakah yang merupakan Kepala Sekolah pertama di SMA Negeri Situraja dan kapan beliau menjabat?',
]

class MyUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def answer_question(self):
        headers = {
            'Content-Type': 'application/json',
        }

        randQuestion = test_data[randrange(len(test_data))]

        response = self.client.post(f'/api/answer?question={randQuestion}', headers=headers)

        if response.status_code != 200:
            print(f"Error: {response.status_code}, Response Text: {response.text}")
