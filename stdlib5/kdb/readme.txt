stdlib5.kdb : 간단한 설정 파일을 저장할 수 있는 키-값 구조체 DB이다.

<python>
class toolbox
    func readstr(str raw)
    # 문자열을 받아 내부 저장소에 파싱 결과를 저장.
    func readfile(str path)
    # 텍스트 파일 경로를 받아 읽고 내부 저장소에 파싱 결과를 파싱.
    func writestr() -> str
    # 내부 저장소 데이터를 포매팅해 문자열로 출력, 데이터 순서를 보장하지 않음.
    func writefile(str path)
    # 내부 저장소 데이터를 포매팅해 텍스트 파일로 출력, 데이터 순서를 보장하지 않음.
    func writestrs() -> str
    # 내부 저장소 데이터를 포매팅해 문자열로 출력, 데이터 순서를 보장함.
    func writefiles(str path)
    # 내부 저장소 데이터를 포매팅해 텍스트 파일로 출력, 데이터 순서를 보장함.
    func getdata(str name) -> [int index, int type, int ptr, any data]
    # 키 문자열을 받아 파라미터와 값으로 이루어진 파이썬 리스트 반환.
    func fixdata(str name, any data)
    # 키 문자열과 새 값을 받아 기존 값을 교체.
    func imp([str name, any data, str end][] arr)
    # 분해된 kdb 리스트를 받아 내부 저장소에 추가, 분해 리스트는 표준 형식이여야 함. (허용 형식은 오류 발생 가능)
    func exp() -> [str name, any data, str end][]
    # 내부 저장소를 분해 kdb 리스트로 반환, 표준 형식임. (데이터 순서를 보장하지 않음)

    .name (int index)[str name] # 키 문자열 -> 인덱스 정수의 딕셔너리.
    .type int[] # 타입 정보를 저장하는 리스트.
    .ptr int[] # 데이터 포인터를 저장하는 리스트.
    .fmem float[] # 실수 값을 저장하는 리스트.
    .cmem complex[] # 복소수 값을 저장하는 리스트.
    .bmem bytes[] # 바이트 값을 저장하는 리스트.

<go>
func Set(interface v) -> kdbvar
# 7type을 받고 그 값을 가지는 kdbvar 반환.
struct kdbvar
    .Dat0 str # 타입 문자열. (nah, bool, int, float, complex, bytes, str)
    .Dat1 bool # 불 타입 저장소.
    .Dat2 int # 정수 타입 저장소.
    .Dat3 float # 실수 타입 저장소. (float64)
    .Dat4 complex # 복소수 타입 저장소. (complex128)
    .Dat5 byte[] # 바이트 배열 저장소.
    .Dat6 string # 문자열 값 저장소.

func Init() -> toolbox
# 내부 map을 초기화, toolbox 구조체를 반환.
struct toolbox
    func Readstr(str* raw)
    # 문자열을 받아 내부 저장소에 파싱 결과를 저장.
    func Readfile(str path)
    # 텍스트 파일 경로를 받아 읽고 내부 저장소에 파싱 결과를 파싱.
    func Writestr() -> *str
    # 내부 저장소 데이터를 포매팅해 문자열로 출력, 데이터 순서를 보장하지 않음.
    func Writefile(str path)
    # 내부 저장소 데이터를 포매팅해 텍스트 파일로 출력, 데이터 순서를 보장하지 않음.
    func Writestrs() -> *str
    # 내부 저장소 데이터를 포매팅해 문자열로 출력, 데이터 순서를 보장함.
    func Writefiles(str path)
    # 내부 저장소 데이터를 포매팅해 텍스트 파일로 출력, 데이터 순서를 보장함.
    func Getpara(str* name) -> [int index, int type, int ptr]
    # 키 문자열을 받아 파라미터로 이루어진 정수 배열 반환.
    func Getvalue(int type, int ptr) -> *kdbvar
    # 파라미터를 받아 값 kdbvar 반환.
    func Getdata(str* name) -> *kdbvar
    # 키 문자열을 받아 값 kdbvar를 반환.
    func Fixdata(str* name, interface v)
    # 키 문자열과 새 값을 받아 기존 값을 교체.
    func Imp(str[]* names, kdbvar[]* datas, str[] ends)
    # 분해된 kdb 배열을 받아 내부 저장소에 추가, 분해 배열은 표준 형식이여야 함. (허용 형식은 오류 발생 가능)
    func Exp() -> str[]* names, kdbvar[]* datas, str[] ends
    # 내부 저장소를 분해 kdb 배열로 반환, 표준 형식임. (데이터 순서를 보장함)

    .Name (int index)[str name] # 키 문자열 -> 정수 인덱스로의 해시맵.
    .Tp byte[] # 타입 정보를 저장하는 바이트 배열.
    .Ptr int[] # 데이터 포인터를 저장하는 정수 배열.
    .Fmem float[] # 실수 값을 저장하는 배열.
    .Cmem complex[] # 복소수 값을 저장하는 배열.
    .Bmem byte[][] # 바이트 배열 값을 저장하는 배열.

효율적인 저장공간 사용과 이식성을 위해 kdb는 내부적으로 키 -> 인덱스 -> 타입, 포인터 -> 값 구조를 가진다.
타입은 0~255 정수이며, 16으로 나눈 몫이 0이면 end가 줄바꿈, 1이면 end가 세미콜론이다.
16으로 나눈 나머지가 0~6이면 각각 nah, bool, int, float, complex, bytes, str이다.
타입과 포인터 배열은 항상 전체 데이터 개수와 동일하며, 인덱스를 통해 접근할 수 있다.
타입 배열의 순서가 실제 데이터의 입력 순서이다.
포인터 정수는 값을 직접 나타내거나, 타입에 따라 다른 값 저장 배열을 가리킬 수 있다.
포인터 정수의 특징 : nah면 항상 0이다, bool이면 참이면 1 거짓이면 0이다.
int면 값 그 자체다, float면 fmem 인덱스다, complex면 cmem 인덱스다.
bytes면 bmem 인덱스다, str이면 utf-8인코딩한 바이트열이 위치한 bmem의 인덱스다.

기본적으로 문자열 키 - 다양한 타입의 값 구조이며, 가능한 값 종류와 문자열 표기 규칙은 다음과 같다.
NULL - nah 로 표기한다.
BOOL - True 와 False 로 표기한다.
INT - 부호와 0~9 숫자. ex) 0, -3, 006
FLOAT - 부호와 0~9 숫자에 소수점이 있는 경우. ex) 0.0, -3.05, 7.0
COMPLEX - 두 숫자형과 부호에 끝이 i인 경우. ex) 0.0+0i, -4-6i, 0.04+4.06i
BYTES - 단일따옴표 안에 16진법 숫자로 표시. ex) '', '00414f', '8A91'
(16진 표기에 소문자가 원칙이나 대문자도 허용한다.)
STRING - 쌍따옴표 안에 #을 이스케이프 코드로 표시. ex) "", "안녕abc01", "31.07"
(#은 ##, "은 #", 줄바꿈은 #n, 공백은 #s가 원칙이나 그냥 사용하는 것도 허용한다.)
("\n"과 "#n"은 길이가 2와 1로 다르다.) ex) "st ri#sng", "#n#####"", "#"#"", "#s\n#n##"

데이터는 다음 구조가 반복된다.
-> 글자 제한이 있는 식별자, 등호 할당자, 문자열로 포메팅된 값, 종결자
식별자엔 공백과 등호를 사용할 수 없다. ex) 가나, wow65, 01
등호 전후로 식별자와 값이 구분된다.
포메팅 방식은 위에서 설명한 방식 그대로다.
종결자는 줄바꿈 혹은 세미콜론으로 구문을 끝낸다.
다음은 데이터 구조 예시이다.
01=01; 가나abc3 = 3
!!= " var = 0 " ; a=6.5;b=6.6
*hey? = ";";:=True

식별자는 마침점을 통해 구조화될 수 있다.
abc.def.ghi와 같은 이름은 이 식별자가 상위 2개의 구조체 안에 있음을 나타낸다.
원칙적으로는 상위 식별자가 모두 텍스트 안에 있을 필요는 없으며,
상위 식별자를 생략하지 않고 모두 명시해야 한다.
다음은 데이터 텍스트 예시이다.
f0.g0.h0 = 0
f0.g1 = 1
f1.g0 = 2
f1 = 3
f0 = 4

가장 앞 부분에서 연속하는 마침점을 사용해서 상위 구조를 생략하는 것도 허용된다.
또한 마침점을 사용하는 것이 원칙이나, 슬래시를 사용하는 것도 허용된다.
등호로 할당되지 않은 식별자는 주석 처리되어 무시된다.
다음은 텍스트 예시이다.
f0 = 0
.g0 = 1
..h0 = 2
.g1 = 3
..h0 = 4
축약을 사용했으므로 상위 구조가 앞서서 모두 등장해야 한다.
com?"#
f1 = 5
f2 = 6;/g0=7;//h0=8;///k=9
