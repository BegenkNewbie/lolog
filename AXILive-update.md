# 🚀 AXILive Update Summary: Rendering Type & Crossfade Dinamis

Dokumen ini merangkum seluruh pembaruan sistem yang telah kita lakukan, mencakup perbaikan *database*, integrasi spreadsheet, hingga fitur *rendering* video tingkat lanjut menggunakan FFmpeg.

## 🛠️ Ringkasan Perubahan

### 1. Penambahan `render_type`
- Menambahkan kolom `render_type` pada *database* SQLite di tabel `live_schedule`.
- Memperbaiki kode migrasi (*Auto Alter Table*) agar inisialisasi nilai *default* kolom berjalan dengan sempurna.
- Mengubah *Struct* `LiveSchedule` di backend Go dan menyertakannya pada semua operasi CRUD (Insert, Read, Update, Delete) di `repo/live.go`.

### 2. Penyesuaian Kolom Spreadsheet
Urutan pembacaan kolom dari *Google Spreadsheet* (dari A ke T / 20 kolom) telah disesuaikan dengan format terbaru Anda.
Urutan khusus yang diperbarui (dimulai dari index 9):
`Tags` ➔ `Privacy` ➔ `Video Intro` ➔ `Video Intro Type` ➔ `Video CTA` ➔ `Rendering Type` ➔ `Altered Content (yes/no)`

### 3. Crossfade Dinamis (*Music Video*)
Sistem penggabungan video musik (*Music Type*) sekarang jauh lebih pintar!
- **Parsing Pintar**: *Backend* otomatis membaca nilai teks dari kolom `Rendering Type` dan mengambil angka durasinya secara dinamis.
- **Kalkulasi Overlap**: Perhitungan batas maksimal panjang video (target) sekarang memperhitungkan total durasi *crossfade* yang terpotong secara akurat.
- **FFmpeg Filter**: Injeksi variabel durasi ke dalam argumen FFmpeg `xfade` (transisi).

### 4. Perbaikan Kompatibilitas Video (QuickTime/Mac)
> [!TIP]
> **Mengapa sebelumnya video crossfade tidak bisa dibuka?**
> Filter kompleks di FFmpeg seperti `xfade` terkadang mengeluarkan *pixel format* mentah yang tidak ramah di semua *player*. 
- **Solusi**: Menyuntikkan perintah `-pix_fmt yuv420p` ke dalam mesin `renderCrossfade` untuk menjamin hasil render 100% kompatibel di Windows, Mac (QuickTime), maupun HP.
- **Log Pintar**: Sistem *error logging* diperbarui. Jika FFmpeg gagal memproses efek, pesan kegagalannya (bukan sekadar exit status 1) akan dicetak dan *file* rusak otomatis dihapus.

### 5. Penyelamatan *Render ASMR*
> [!WARNING]
> Sebelumnya, jika video tipe ASMR hanya berisi **1 video polos** (tanpa diloop dan tanpa intro), proses standarisasi resolusi terlewatkan.
- **Solusi**: Menambahkan kondisional khusus di `RenderVideoASMR`. Jika sistem memutuskan untuk *skip* `concatVideos`, video utama akan langsung diarahkan ke `normalizeVideo` agar format resolusinya seragam dan *watermark* berhasil menempel.

---

## 📋 Aturan Main (Rules of the Game)

Berikut adalah panduan dan format yang harus digunakan tim Anda saat mengisi data di *Google Spreadsheet*.

### Pengisian Kolom `Rendering Type`
Kolom ini menentukan jenis penyambungan antar video (terutama untuk tipe *Music*).

| Teks di Spreadsheet | Hasil / Efek yang Terjadi |
| :--- | :--- |
| **`concat`** *(atau kosong)* | Penyambungan video standar tanpa efek (*cut* biasa). Sangat cepat dan irit CPU. |
| **`crossfade`** | Efek transisi halus antar video. Durasi *default* adalah **1.0 detik**. |
| **`crossfade-1.5`** | Transisi *crossfade* berdurasi **1.5 detik**. |
| **`crossfade-2`** | Transisi *crossfade* berdurasi **2.0 detik**. |
| **`crossfade-3.5`** | Transisi *crossfade* berdurasi **3.5 detik**. |

> [!NOTE]
> - Sistem membaca dengan format **Case-Insensitive** (huruf besar/kecil seperti `Crossfade-2` atau `CROSSFADE` tetap akan terbaca normal).
> - Jika ada kesalahan pengetikan angka (contoh: `crossfade-abc`), sistem secara otomatis akan menyelamatkannya dengan mengubahnya menjadi durasi *default* `1.0` detik, sehingga tidak akan membuat server *crash*.

### Pengisian Kolom `Video Intro Type`
Kolom ini menentukan bagaimana sistem akan memproses video intro yang disisipkan (apakah dipakai apa adanya, atau butuh proses *chromakey/green-screen* dengan video latar).

| Teks di Spreadsheet | Hasil / Efek yang Terjadi |
| :--- | :--- |
| **`fixed`** | Video intro akan di-download dan diputar **apa adanya** di depan video utama. Tidak ada proses *chromakey* atau penggabungan latar belakang. |
| **(Kosong)** | *Default*. Sistem berasumsi video intro adalah *green screen*. Sistem akan otomatis mengambil **2 video pertama** (*Music*) atau **1 video utama** (*ASMR*), menggabungkannya, lalu menempelkannya di latar belakang (*chromakey*). |
| **`1`, `2`, `3`, dst** | *(Khusus Music)* Menentukan **berapa banyak video latar** yang harus disambung terlebih dahulu sebelum dijadikan latar belakang intro *green screen*. Contoh: `3` berarti 3 video latar pertama akan digabung dulu untuk menemani durasi intro. |

> [!TIP]
> Jika kolom ini diisi dengan teks sembarangan selain `fixed` atau angka (misal diisi teks `random`), sistem tidak akan *crash* dan akan menganggapnya kosong (otomatis menggunakan *default* 2 video latar).

### Flow Pengecekan Video
Jika ada video yang dilaporkan gagal di-render, pastikan:
1. Resolusi video *source* awal semuanya memiliki FPS yang masuk akal (hindari mencampur video 24fps dan 60fps dengan ekstrem jika memungkinkan, walau `normalizeVideo` berusaha menyamakannya).
2. Periksa log terminal jika FFmpeg gagal. Kini log FFmpeg memuat *error stream* yang sangat mendetail.
