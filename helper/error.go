package helper

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}

	/*
		NOTE!!
			- untuk troubleshooting/ pesan errornya jika menggunakan panic(err) bisa melihat di console apa yang menjadi panic errornya
			- perhatikan, bahwa lebih baik melakukan error handler dengan memberikan panic, ditakutkan apbila hanya melakukan println errornya apa,
			- validasinya tetap akan masuk ke database, oleh karena itu dengan adanya panic langsung mencegah hal ini dan otomatis akan menolaknya
	*/
}
