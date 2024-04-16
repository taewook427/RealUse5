stdlib5.kio [win & linux] : 터미널/파일 입출력 공통 함수들.

<go>
func Print(any v, float t)
# t초 대기 후 v 내용 출력.
func Input(str q) -> str
# 질문 q 출력 후 \r\n 제거한 텍스트 응답 반환.
func Bequal(byte[] a, byte[] b) -> bool
# 두 바이트 슬라이스가 동일한지 확인.
func Bread(str raw) -> (byte[], error)
# 문자열 표현된 바이트를 바이트 슬라이스로 변환. (ex. "8a9C" -> B)
func Bprint(byte[] raw) -> str
# 바이트 슬라이스를 문자열로 변환. (ex. B -> "0c17")
func Abs(str path) -> str
# 입력 경로를 표준경로로 변환. (슬래시 + 폴더 슬래시 종결 + 절대 경로)
func Size(str path) -> int
# 파일 크기를 반환. 존재하지 않는 파일인 경우 -1 반환.
func Open(str path, str mode) -> (OSFILE*, error)
# 파일 경로와 모드를 받아 운영체제 입출력 파일 포인터 반환.
"r" : 읽기 전용, "w" : 덮어쓰기, "a" : 이어쓰기, "x" : 실행권한 있는 덮어쓰기.
func Read(OSFILE* f, int size) -> (byte[], error)
# 일정 크기만큼 파일 읽기. size가 음수라면 전체 읽기. 읽기 성공한 크기만을 자동으로 잘라 반환.
func Write(OSFILE* f, byte[] data) -> (byte[], error)
# 파일 쓰기. 쓰기 실패한 바이트만큼 자동으로 잘려 반환됨.

!!!!! Reading 경고 !!!!!
Golang의 기본 읽기 함수(Read)는 1GiB를 초과하는 큰 파일을 한번에 제대로 읽을 수 없습니다.
파일 포인터에 직접 Read를 사용할 경우, 1GiB 이하의 단위로 끊어 읽으십시오.
kio.Read를 사용할 경우, 속도는 더 느리나 자동으로 큰 파일을 끊어 읽어줍니다.

많은 코드에서 공통적으로 쓰일 만한 기능들을 포장했습니다.
