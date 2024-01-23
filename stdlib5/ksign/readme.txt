stdlib5.ksign : RSA-pss 전자서명과 폴더-파일 통합 해싱.

<python.st>
func khash(str path) -> bytes hash
# 절대/상대 윈도우/리눅스 파일/폴더 경로를 받고 해싱 결과 64B 반환.
func genkey(int n) -> str public, str private
# n바이트 강도의 RSA 키 생성. RSA2048 : 256, RSA 4096 : 512
func sign(str private, bytes plain) -> bytes enc
# 80B plain 평문을 개인키로 서명. 암호화 결과물 반환.
func verify(str public, bytes enc, bytes plain) -> bool isvalid
# 공개키, 암호화 바이트, 80B plain을 받아 올바른 서명인지 검증.
func fm(str name, bytes hashed) -> bytes plain
# 서명 이름 문자열, 원본 메세지 해시 64B를 받아 plain 80B로 패키징.

<python.hy>
class toolbox
    func khash(str path) -> bytes hash
    # 절대/상대 윈도우/리눅스 파일/폴더 경로를 받고 해싱 결과 64B 반환.
    func genkey(int n) -> str public, str private
    # n바이트 강도의 RSA 키 생성. RSA2048 : 256, RSA 4096 : 512
    func sign(str private, bytes plain) -> bytes enc
    # 80B plain 평문을 개인키로 서명. 암호화 결과물 반환.
    func verify(str public, bytes enc, bytes plain) -> bool isvalid
    # 공개키, 암호화 바이트, 80B plain을 받아 올바른 서명인지 검증.
    func fm(str name, bytes hashed) -> bytes plain
    # 서명 이름 문자열, 원본 메세지 해시 64B를 받아 plain 80B로 패키징.

<go.st>
func Khash(str path) -> byte[] hash
# 절대/상대 윈도우/리눅스 파일/폴더 경로를 받고 해싱 결과 64B 반환.
func Genkey(int n) -> str public, str private
# n바이트 강도의 RSA 키 생성. RSA2048 : 256, RSA 4096 : 512
func Sign(str private, byte[] plain) -> byte[] enc
# 80B plain 평문을 개인키로 서명. 암호화 결과물 반환.
func Verify(str public, byte[] enc, byte[] plain) -> bool isvalid
# 공개키, 암호화 바이트, 80B plain을 받아 올바른 서명인지 검증.
func Fm(str name, byte[] hashed) -> byte[] plain
# 서명 이름 문자열, 원본 메세지 해시 64B를 받아 plain 80B로 패키징.

!!! 중요 경고 !!!
- 이름 문자열은 길이가 16바이트를 넘을 수 없습니다.
- plain은 80바이트, hashed는 64바이트 길이입니다.
- standard 버전은 py와 go가 호환되지 않습니다.
khash는 대규모 폴더에서 값이 달라짐, 전자서명은 호환 불가.
go standard와 호환되는 hybrid 버전을 사용하십시오.
- 키 생성 시 n 비트가 아니라 n 바이트 단위입니다.
- khash를 너무 크기가 큰 파일에 대해 사용하지 마십시오.

!!! python VM 경고 !!!
threading 모듈 사용 시 메인 스크립트에
if __name__ == '__main__':
문구를 추가할 것.

khash는 다음 규칙을 따라 64바이트 결과를 내보냅니다.
1. 빈 파일의 해시값은 각 위치가 모두 0인 길이 64의 바이트.
2. 빈 폴더의 해시값은 빈 파일과 동일하게 계산.
3. 비지 않은 파일의 해시값은 파일 바이너리를 SHA3-512 해싱한 결과물.
4. 비지 않은 폴더의 해시값은 다음과 같이 계산함.
    4-1. 현재 폴더의 위치를 정규화 (리눅스 현식, 뒤가 /로 끝나게)
    4-2. 현재 폴더의 하위 항목의 이름에 파일은 그대로, 폴더는 뒤에 /을 붙임.
    4-3. 파일/폴더 상관없이 오름차순 정렬함.
    4-4. 각 항목의 해시값 64바이트를 정렬된 순으로 이어붙여 하나의 바이너리를 만듦.
    4-5. 폴더의 해시값은 이 바이너리를 SHA3-512 해싱한 결과물.
