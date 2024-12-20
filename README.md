# Jabar Coding Camp Partnership - Project Challenge Batch 2
![Logo JCC](https://scontent.fsub3-2.fna.fbcdn.net/v/t39.30808-6/216740704_101320238899855_3589530600193911774_n.png?_nc_cat=104&ccb=1-7&_nc_sid=5f2048&_nc_eui2=AeFr9AMgwhamgYyESSGlyt3TsqpN-PPBq86yqk3488GrzjnTvhM1BtaOwWwgPv-1_rmO5vHfGs1_p3wfD1mIigJU&_nc_ohc=lVnx0FX2JHIAX89q3_W&_nc_zt=23&_nc_ht=scontent.fsub3-2.fna&oh=00_AfBetPTDlsRN2n8p_krAKj8I6gLJ3knb0sNJP1AOmjdxVw&oe=65553922)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)

### Informasi Peserta  

Nama: 
Kelas: Backend Golang 


### Link  

~~[Link API di Heroku (swagger)](https://get-nearby-places-jcc.herokuapp.com/swagger/index.html)~~   (sudah tidak aktif seiring free tier heroku sudah berhenti)  

[Link YouTube Demo Penggunaan API](https://youtu.be/-zilA1NbZS8)

## Deskripsi API
Web service (API) untuk melakukan pencarian daftar tempat/place terdekat. Tempat/place digenerate secara random, memiliki atribut dan constraint yang melekat pada region (kota/kabupaten/desa/keluarahan/kecamatan) yang ditempatinya. API akan mengembalikan semua tempat terdekat dalam radius 5 km dari sebuah titik koordinat latitude dan longitude yang diberikan.

## Penggunaan

Install seluruh package yang diperlukan
```bash
go install
```

Jalankan program
```bash
go run main.go
```

Kunjungi halaman swagger di
```bash
http://localhost:8080/
```

## NOTE   

Kode telah menggunakan fitur generic yang mulai ada pada golang versi 1.18. Selengkapnya tentang generic: 

[Release Note Go 1.18](https://tip.golang.org/doc/go1.18#generics)   


Pastikan versi golang anda >= 1.18
## Disclamier  

Tidak diizinkan menyalin/mengutip kode penuh/sebagian untuk apapun tugas JCC next batch atau batch saat ini (2) di kelas apapun.  
Boleh digunakan hanya untuk referensi ide.
