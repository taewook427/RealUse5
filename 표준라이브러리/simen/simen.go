package simen // test538 : stdlib khash/simen

import (
	"bytes"
)

type toolbox struct {
	sbox    [256]byte
	invsbox [256]byte
	bit     [256][8]bool
	key     [][8][8]byte
}

// 시작 전 기본 데이터 생성
func Init() toolbox {
	output := toolbox{}
	prebox := [16]byte{14, 4, 13, 1, 2, 15, 10, 6, 8, 3, 11, 9, 5, 12, 7, 0}
	for i := 0; i < 256; i++ {
		output.sbox[i] = 16*prebox[i/16] + prebox[i%16]
		num := i
		div := 128
		for j := 0; j < 8; j++ {
			if num < div {
				output.bit[i][j] = false
			} else {
				num = num - div
				output.bit[i][j] = true
			}
			div = div / 2
		}
	}
	for i := 0; i < 256; i++ {
		output.invsbox[output.sbox[i]] = byte(i)
	}
	return output
}

func shift(array *[16]byte, value byte) { // shift 밀기 연산, 0~31
	add := value % 4
	ptr := value / 4
	if ptr < 4 {
		temp := [8]byte{array[4*ptr], array[4*ptr+1], array[4*ptr+2], array[4*ptr+3], array[4*ptr], array[4*ptr+1], array[4*ptr+2], array[4*ptr+3]}
		array[4*ptr] = temp[add]
		array[4*ptr+1] = temp[add+1]
		array[4*ptr+2] = temp[add+2]
		array[4*ptr+3] = temp[add+3]
	} else {
		ptr = ptr - 4
		temp := [8]byte{array[ptr], array[ptr+4], array[ptr+8], array[ptr+12], array[ptr], array[ptr+4], array[ptr+8], array[ptr+12]}
		array[ptr] = temp[add]
		array[ptr+4] = temp[add+1]
		array[ptr+8] = temp[add+2]
		array[ptr+12] = temp[add+3]
	}
}

func invshift(array *[16]byte, value byte) { // invshift 밀기 연산 역함수
	add := value % 4
	ptr := value / 4
	shift(array, 4*ptr+(4-add)%4)
}

func logic(mode int, x *[8]bool, y *[8]bool) byte { // logic 논리 합성 연산
	output := [8]bool{}
	if mode == 0 {
		for i := 0; i < 8; i++ {
			output[i] = (!x[i] && y[i]) && (x[i] || y[i])
		}
	} else if mode == 1 {
		for i := 0; i < 8; i++ {
			output[i] = (x[i] && !y[i]) || (x[i] != y[i])
		}
	} else if mode == 2 {
		for i := 0; i < 8; i++ {
			output[i] = (!x[i] || y[i]) || (x[i] && y[i])
		}
	} else {
		for i := 0; i < 8; i++ {
			output[i] = (x[i] != y[i]) && (x[i] || !y[i])
		}
	}
	var temp byte = 0
	for i := 0; i < 8; i++ {
		if output[i] {
			temp = temp + 1<<(7-i)
		}
	}
	return temp
}

func xor(x *[8]bool, y *[8]bool) byte { // xor 숫자 출력
	output := [8]bool{}
	for i := 0; i < 8; i++ {
		output[i] = (x[i] != y[i])
	}
	var temp byte = 0
	for i := 0; i < 8; i++ {
		if output[i] {
			temp = temp + 1<<(7-i)
		}
	}
	return temp
}

func (db *toolbox) revhash(x *[8]byte, n int) [8]byte { // 뒤집고 해싱, n칸 밀기
	temp := make([]byte, 8)
	for i := 0; i < 8; i++ {
		temp[i] = x[7-i]
	}
	t := db.Hash(temp)
	output := [8]byte{}
	for i := 0; i < 8; i++ {
		output[i] = t[(i+n)%8]
	}
	return output
}

// 해시 함수 input [n]byte -> hash [8]byte
func (db *toolbox) Hash(data []byte) [8]byte {
	// 8바이트 패딩
	length := len(data)
	padlen := (8 - length%8) % 8
	for i := 0; i < padlen; i++ {
		data = append(data, byte(padlen))
	}
	length = length + padlen
	chunk := make([][8]byte, length/8)
	for i := 0; i < length/8; i++ {
		for j := 0; j < 8; j++ {
			chunk[i][j] = data[8*i+j]
		}
	}

	iv := [8]byte{66, 75, 235, 179, 145, 234, 180, 128}
	for _, r := range chunk {
		for i := 0; i < 8; i++ {
			iv[i] = xor(&db.bit[iv[i]], &db.bit[r[i]]) // iv ^ chunk
		}
		matrix := [16]byte{}
		for i := 0; i < 4; i++ { // logic operation
			matrix[4*i] = logic(0, &db.bit[iv[i]], &db.bit[iv[7-i]])
			matrix[4*i+1] = logic(1, &db.bit[iv[i]], &db.bit[iv[7-i]])
			matrix[4*i+2] = logic(2, &db.bit[iv[i]], &db.bit[iv[7-i]])
			matrix[4*i+3] = logic(3, &db.bit[iv[i]], &db.bit[iv[7-i]])
		}
		for _, l := range r { // shift operation
			shift(&matrix, l%32)
			shift(&matrix, l/8)
		}
		for i, l := range matrix { // sbox 치환
			matrix[i] = db.sbox[l]
		}
		// matrix 합치기
		iv[0] = matrix[14] & matrix[4]
		iv[1] = xor(&db.bit[matrix[13]], &db.bit[matrix[1]])
		iv[2] = xor(&db.bit[matrix[2]], &db.bit[matrix[15]])
		iv[3] = matrix[10] & matrix[6]
		iv[4] = xor(&db.bit[matrix[8]], &db.bit[matrix[3]])
		iv[5] = matrix[11] | matrix[9]
		iv[6] = matrix[5] | matrix[12]
		iv[7] = xor(&db.bit[matrix[7]], &db.bit[matrix[0]])
	}

	return iv
}

// 키 설정 []byte
func (db *toolbox) Setkey(data []byte) {
	// 8바이트 패딩
	length := len(data)
	if length == 0 {
		panic("keylength0error")
	}
	padlen := (8 - length%8) % 8
	for i := 0; i < padlen; i++ {
		data = append(data, byte(padlen))
	}
	length = length + padlen
	chunk := make([][8]byte, length/8)
	for i := 0; i < length/8; i++ {
		for j := 0; j < 8; j++ {
			chunk[i][j] = data[8*i+j]
		}
	}

	db.key = make([][8][8]byte, len(chunk))
	for i, r := range chunk {
		k0a := r
		k0b := db.revhash(&k0a, 5)

		k1a := [8]byte{}
		for j := range k0a {
			k1a[j] = xor(&db.bit[k0a[j]], &db.bit[k0b[j]])
		}
		k1b := db.revhash(&k1a, 3)

		k2a := [8]byte{}
		for j := range k1a {
			k2a[j] = xor(&db.bit[k1a[j]], &db.bit[k1b[j]])
		}
		k2b := db.revhash(&k2a, 7)

		k3a := [8]byte{}
		for j := range k2a {
			k3a[j] = xor(&db.bit[k2a[j]], &db.bit[k2b[j]])
		}
		k3b := db.revhash(&k3a, 1)

		k4a := [8]byte{}
		for j := range k3a {
			k4a[j] = xor(&db.bit[k3a[j]], &db.bit[k3b[j]])
		}
		k4b := db.revhash(&k4a, 5)

		db.key[i] = [8][8]byte{k1a, k1b, k2a, k2b, k3a, k3b, k4a, k4b}
	}
}

// 암호화
func (db *toolbox) Encrypt(data []byte) []byte {
	// 16바이트 패딩
	length := len(data)
	padlen := 16 - length%16
	for i := 0; i < padlen; i++ {
		data = append(data, byte(padlen))
	}
	length = length + padlen
	chunk := make([][16]byte, length/16)
	for i := 0; i < length/16; i++ {
		for j := 0; j < 16; j++ {
			chunk[i][j] = data[16*i+j]
		}
	}

	output := []byte{}
	writer := bytes.NewBuffer(output)
	for i, r := range chunk {
		matrix := r
		key := db.key[i%len(db.key)]

		for j := 0; j < 4; j++ { // 4 rounds

			for k, u := range matrix { // xor
				if k%4 < 2 {
					matrix[k] = xor(&db.bit[u], &db.bit[key[2*j][k/2+k%2]])
				} else {
					matrix[k] = xor(&db.bit[u], &db.bit[key[2*j+1][k/2+k%2-1]])
				}
			}

			for k := 0; k < 8; k++ { // 8 shift
				shift(&matrix, (key[2*j][k]+key[2*j+1][k])%32)
			}

			for k, u := range matrix {
				matrix[k] = db.sbox[u] // sbox 치환
			}

		}

		for _, l := range matrix {
			writer.WriteByte(l)
		}
	}

	return writer.Bytes()
}

func (db *toolbox) Decrypt(data []byte) []byte {
	length := len(data)
	num := length / 16
	chunk := make([][16]byte, num)
	for i := 0; i < num; i++ {
		for j := 0; j < 16; j++ {
			chunk[i][j] = data[16*i+j]
		}
	}
	output := []byte{}
	writer := bytes.NewBuffer(output)

	for i, r := range chunk {
		matrix := r
		key := db.key[i%len(db.key)]

		for jj := 0; jj < 4; jj++ { // 4 rounds
			j := 3 - jj

			for k, u := range matrix {
				matrix[k] = db.invsbox[u] // sbox 치환
			}

			for kk := 0; kk < 8; kk++ { // 8 shift
				k := 7 - kk
				invshift(&matrix, (key[2*j][k]+key[2*j+1][k])%32)
			}

			for k, u := range matrix { // xor
				if k%4 < 2 {
					matrix[k] = xor(&db.bit[u], &db.bit[key[2*j][k/2+k%2]])
				} else {
					matrix[k] = xor(&db.bit[u], &db.bit[key[2*j+1][k/2+k%2-1]])
				}
			}
		}

		for _, l := range matrix {
			writer.WriteByte(l)
		}
	}

	output = writer.Bytes()
	padlen := output[len(output)-1]
	return output[0 : len(output)-int(padlen)]
}
