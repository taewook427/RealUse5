stdlib5.picdt [win & linux] : gen5 project에 공통적으로 쓰이는 사진 파일 데이터.

<python>
class toolbox
    .data0 # kzip5 png 128x128 도안
    .data1 # kzip5 webp 128x128 도안
    .data2 # kpic5 png 128x128 도안
    .data3 # kpic5 webp 128x128 도안
    .data4 # kaes5 png 128x128 도안
    .data5 # kaes5 webp 128x128 도안
    .data6 # zipre5 png 128x128 도안
    .data7 # zipre5 webp 128x128 도안

<go>
func Kz5png() -> byte[]*
# kzip5 png 128x128 도안
func Kz5webp() -> byte[]*
# kzip5 webp 128x128 도안
func Kp5png() -> byte[]*
# kpic5 png 128x128 도안
func Kp5webp() -> byte[]*
# kpic5 webp 128x128 도안
func Ka5png() -> byte[]*
# kaes5 png 128x128 도안
func Ka5webp() -> byte[]*
# kaes5 webp 128x128 도안
func Zr5png() -> byte[]*
# zipre5 png 128x128 도안
func Zr5webp() -> byte[]*
# zipre5 webp 128x128 도안

소스코드 바이너리 적재기를 통해 생성된 코드를 수정하였습니다.
