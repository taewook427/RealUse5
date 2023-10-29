stdlib5.simen : 단순 암호화, 해시 기능

<python>
class toolbox
    func hash(bytes input) -> bytes output
    # 임의 길이 바이트를 입력받아 8바이트 길이의 해시값을 출력
    func setkey(bytes key)
    # 0이 아닌 길이의 바이트를 입력받아 키로 설정
    func encrypt(bytes input) -> bytes output
    # 임의 길이 바이트를 입력받아 설정된 키로 암호화 후 반환
    func decrypt(bytes input) -> bytes output
    # 16n 길이 바이트를 입력받아 설정된 키로 복호화 후 반환

<go>
func Init() -> toolbox module
# toolbox 구조체를 초기화한 후 반환하는 초기화 함수
struct toolbox
    func Hash(byte[] input) -> byte[8] output
    # 임의 길이 바이트를 입력받아 8바이트 길이의 해시값을 출력
    func Setkey(byte[] key)
    # 0이 아닌 길이의 바이트를 입력받아 키로 설정
    func Encrypt(byte[] input) -> byte[] output
    # 임의 길이 바이트를 입력받아 설정된 키로 암호화 후 반환
    func Decrypt(byte[] input) -> byte[] output
    # 16n 길이 바이트를 입력받아 설정된 키로 복호화 후 반환

추가 라이브러리 없이 간단하게 해싱, 암호화를 할 수 있는 모듈입니다.
sbox 치환 연산 (1 바이트를 비선형적으로 1 바이트로 치환)
shift 연산 (4x4 바이트 행렬의 각 행or열 0~3칸씩 밀기)
logic 연산 (and, or, xor을 복합한 논리 연산)
revhash 연산 (8 바이트 데이터를 뒤집고 해싱 후 n칸 밀기)

해싱은 다음 알고리즘으로 진행됩니다.
1. 8바이트 패딩 및 나누기
# 데이터가 8의 배수 길이가 아니라면 패딩하고 8바이트 단위로 나눕니다.
2. iv에 각 8바이트 데이터 순차 적용
# iv 8바이트에 데이터 8바이트를 적용합니다.
    2-1. iv ^ chunk
    # iv와 chunk의 숫자를 xor하여 iv를 다시 작성합니다.
    2-2. logic(iv) -> matrix
    # iv의 각 항목쌍에 logic 연산을 적용시켜 4x4 matrix를 생성합니다.
    2-3. shift(matrix, chunk % 32), shift(matrix, chunk // 8)
    # chunk의 각 숫자들 만큼 matrix에 shift 연산을 진행합니다.
    2-4. sbox(matrix)
    # matrix의 각 원소를 sbox 치환합니다.
    2-5. matrix -> iv
    # martix를 다시 대응쌍대로 합쳐 다음 라운드의 iv를 생성합니다.
3. return iv
# 모든 chunk에 대해 연산을 거친 iv를 반환합니다.

키 설정은 다음 알고리즘으로 진행됩니다.
1. 8바이트 패딩 및 나누기
# 데이터가 8의 배수 길이가 아니라면 패딩하고 8바이트 단위로 나눕니다.
# 8바이트 라운드 키 8개가 한 세트의 키이며 패딩 후 나뉜 키 청크의 길이만큼의 키 세트가 있습니다.
2. 원본 키 청크 k0a에서 키 세트 k1a ~ k4a, k1b ~ k4b 유도
# k0a를 원본 8 바이트 키로 설정하고 revhash, xor 연산을 통해 길이 8의 키 세트를 유도합니다.
    2-1. revhash(k0a, 5) -> k0b, xor(k0a, k0b) -> k1a
    # k0a에 revhash 연산을 적용해 k0b를 얻고, 둘을 xor 하여 k1a를 얻습니다.
    2-2. revhash(5, 3, 7, 1, 5) -> (k0b, k1b, k2b, k3b, k4b), xor -> (k1a, k2a, k3a, k4a)
    # 순차적으로 revhash, xor 연산을 적용시켜 키 세트 k1a ~ k4b를 얻습니다.
3. class / struct 설정
# 얻어진 키 세트들을 객체나 구조체 변수에 담습니다.

암호화는 다음 알고리즘으로 진행됩니다. 복호화는 암호화를 역순으로 진행합니다.
1. 16바이트 패딩 및 나누기
# 데이터가 16의 배수 길이가 되도록 패딩하고 16바이트 길이로 나눕니다.
2. 데이터 청크별 키 세트 적용
# 각 데이터 청크에 대응하는 키 세트를 적용시켜 암호화합니다.
# 키는 kna/knb이므로 한 청크에서 암호화는 4 라운드 진행됩니다.
    2-1. chunk -> 4x4 matrix, xor(matrix, key) -> matrix
    # chunk를 4x4 matrix로 만들고 라운드 키와 xor 합니다.
    2-2. shift(matrix, key) *8
    # matrix에 (key n a + key n b) % 32만큼 총 8회 shift 연산을 합니다.
    2-3. sbox(matrix) -> matrix
    # matrix를 sbox 치환합니다.
3. matrix 합치기, return bytes
# 연산이 끝난 matrix를 합쳐 bytes로 만들고 반환합니다.
