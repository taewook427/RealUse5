stdlib5.legsup [win & linux] : gen1~4 구세대 데이터 형식 호환 계층.

<go-exp>
common:
    func G3mold() -> byte[]
    # gen3kpic에 사용되는 기본 주형 사진. (png binary)
    func G3zip() -> byte[]
    # gen3kzip에 사용되는 아이콘 사진. (png binary)
    func G3kf() -> byte[]
    # gen3kaes의 기본키파일.
    func G3pic() -> byte[]
    # gen3kaes에 사용되는 아이콘 사진. (png binary)
    func G4kf() -> byte[]
    # gen4kaes의 기본키파일.
    func G4pic() -> byte[]
    # gen4kaes에 사용되는 아이콘 사진. (png binary)
    func Genkf(str path) -> byte[]
    # 경로의 파일을 모두 읽어 반환, 오류 시 nil 반환.

gen1:
    struct G1enc
        func Init()
        # 구조체 내부 필드를 초기화.
        func Encrypt() -> error
        # 구조체 내부 데이터를 참고하여 암호화.
        func Decrypt() -> error
        # 구조체 내부 데이터를 참고하여 복호화.
        func View() -> error
        # 구조체 내부 데이터를 참고하여 암호화된 파일 정보 해석.

        .Path str # 암호화/복호화/해석 대상 파일 경로.
        .Pw str # 비밀번호.
        .Hint str # 비밀번호 힌트.

gen2:
    struct G2enc
        func Init()
        # 구조체 내부 필드를 초기화.
        func Encrypt() -> (str, error)
        # 구조체 내부 데이터를 참고하여 암호화. 암호화 결과 경로 반환.
        func Decrypt() -> error
        # 구조체 내부 데이터를 참고하여 복호화.
        func View() -> error
        # 구조체 내부 데이터를 참고하여 암호화된 파일 정보 해석.

        .Path str # 암호화/복호화/해석 대상 파일 경로.
        .Pw str # 비밀번호.
        .Hint str # 비밀번호 힌트.
        .Hidename bool # 암호화 시 파일명 숨김 옵션. True일 경우 파일명 숨김.

gen3:
    struct G3data
        func Init()
        # gen3kdb 항목 하나를 이루는 데이터 노드를 초기화.
        func Append(G3data* tgt)
        # 데이터 노드의 가장 끝에 데이터 노드 추가. (리스트 형식 데이터 등에 사용)
        func Length() -> int
        # 현재 노드를 시작으로 끝 노드까지의 길이.
        func Locate(int pos) -> G3data*
        # 현재 노드를 기준으로 pos번째 노드 반환.
        func Print(bool zipstr, bool zipexp) -> str
        # 노드 문자열 출력. zipstr/zipexp는 문자열/표현형식 단축 여부.

        .Next *G3data # 다음 데이터 노드 포인터. 없을 경우 nil.
        .Vtype rune # 데이터 노드의 값 종류. ('i', 'f', 's', 'n')
        .IntV int # 정수 값 저장 위치.
        .FloatV float # 실수 값 저장 위치.
        .StrV str # 문자열 값 저장 위치.

    struct G3node
        func Read(str[] frag, rune[] trait) -> error
        # gen3kdb의 kobj 하나에 해당하는 구문을 파싱 후 저장.
        func Write(int indent, bool zipstr, bool zipexp) -> str
        # indent만큼 공백 들여쓰기, 문자열/표현형식 압축 여부에 따라 텍스트 생성.
        func Locate(str name) -> G3node*
        # 현재 레벨에서 name 이름을 가진 하위 노드 반환. 존재하지 않으면 nil 반환.
        func Revise(G3data* tgt) -> error
        # 이 노드의 데이터를 tgt로 수정.

        .Name str # 이 kobj/node의 이름.
        .Data G3data* # 데이터 노드. 객체 역할인 경우 nil.
        .Child G3node[] # 자식 노드들. 데이터 역할인 경우 nil.

    struct G3kdb
        func Read(str raw) -> error
        # gen3kdb 문자열을 입력받아 해석, 내부 데이터 저장.
        func Write() -> str
        # 내부 데이터를 gen3kdb 문자열로 출력.
        func Locate(str name) -> G3node*
        # "#"으로 구분된 주소를 가진 노드 반환. 존재하지 않으면 nil 반환.

        .Zipstr bool # 문자열 압축 여부. 압축 시 공백/줄바꿈이 #표현으로 바뀜.
        .Zipexp bool # 표현형식 압축 여부. 압축 시 들여쓰기와 텍스트 구조에 공백/줄바꿈이 사라짐.
        .Node G3node* # 내부 노드. 이 노드를 설정하여 임의 데이터 구조 생성 가능.

    struct G3kzip
        func Init()
        # 패키징에 필요한 내부 데이터 설정.
        func Packf(str[] tgt, str path) -> error
        # 파일들 패키징. tgt 항목을 패키징해 path에 생성.
        func Packd(str tgt, str path) -> error
        # 폴더 패키징. tgt 폴더를 패키징해 path에 생성.
        func View(str tgt) -> error
        # 패키징된 파일 정보 해석. CRC32 값 검사.
        func Unpack(str tgt) -> error
        # 패키징된 tgt 파일 풀기. "./temp261/" 폴더에 생성됨.

        .Prehead byte[] # 사진 위장 헤더. (1024nB)
        .Header byte[] # 메인헤더 (18B)
        .Chunkpos int[] # 청크 시작 위치. subhead + data가 한 청크.
        .Subhead byte[][] # 각 청크의 subhead.
        .Winsign bool # 폴더 표시를 위해 백슬래시 사용. (윈도우 전용 구형 코드와 호환성)

    struct G3kaesall
        func Init(int core, int chunk)
        # 전체모드 암호화에 필요한 내부 데이터 설정. (0, 0)으로 기본 모드(32, 128k).
        func Encrypt(str path, str pw, byte[] kf) -> error
        # path 파일을 암호화.
        func Decrypt(str path, str pw, byte[] kf) -> error
        # path 파일을 복호화.
        func View(str path) -> error
        # 암호화된 path 파일의 정보 해석.

        .Prehead byte[] # 사진 위장 헤더. (1024nB)
        .Metadata byte[] # 메타데이터 (18B)
        .Mainhead byte[] # 메인헤더
        .Subhead byte[] # 보조헤더 (현재 사용하지 않음)
        .Hidename bool # 이름 숨기기 여부. True일 경우 무작위 숫자 파일명으로 생성됨.
        .Hint str # 비밀번호 힌트.
        .Respath str # 암호화/복호화 결과 파일 경로.

    struct G3kaesfunc
        func Encrypt(str before, str after, byte[] akey) -> error
        # 32B akey로 파일 암호화. (32, 128k) 모드 고정.
        func Decrypt(str before, str after, byte[] akey) -> error
        # 32B akey로 파일 복호화. (core, chunk)는 암호화 파일에 따름.

        .Metadata byte[] # 메타데이터 (18B), 사진위장과 보조헤더는 사용하지 않음.
        .Mainhead byte[] # 메인헤더

    struct G3kv3
        func Encrypt(str pw, byte[] kf, str path) -> error
        # path 폴더를 func+kv3 모드로 암호화.
        func Decrypt(str pw, byte[] kf, str path) -> error
        # path 폴더를 func+kv3 모드로 복호화.
        func View(str path) -> error
        # 암호화된 폴더의 정보 해석.

        .Hint str # kv3 암호화의 비밀번호 힌트.

    struct G3kpic
        func Init(str path, int row, int col) -> error
        # 주형 사진 설정. 빈 문자열로 기본주형사진 사용.
        # 가로세로 크기는 모두 4의 배수 혹은 -1로 사진크기 그대로 설정. png만 사용 가능.
        func Detect(str path) -> (str, int, error)
        # 폴더 안에 gen3kpic 파일을 감지, 이름과 개수 반환. png 모드만 사용함.
        func Pack(str tgt, str exdir) -> (str, int)
        # tgt 파일을 패키징해 exdir 폴더 내부에 사진들 생성. 사진 이름과 개수 반환.
        func Unpack(str path, str tgtdir, str name, int num)
        # tgtdir 폴더에서 name, num 데이터로 gen3kpic 사진들을 path 파일로 복구.

        .Pcover bool # 사진 위장 사용 여부. True면 2배수 모드, False면 1배수 모드.

    func G3picre(byte[] pic, str[] files, str path) -> error
    # files 파일들을 pic 위장사진헤더를 가진 zip으로 path 경로에 압축시킴.

gen4:
    struct G4enc
        func Encrypt(str[] files, byte[] pw) -> (str, error)
        # 파일들 암호화, 첫 번째 암호화 대상 파일과 같은 폴더에 암호화된 파일 생성, 해당 경로 반환됨.
        func Decrypt(str path, byte[] pw) -> error
        # 암호파일과 같은 폴더에 원본 파일 복구.
        func View(str path) -> error
        # 암호화된 파일의 정보 해석.

        .Hint str # gen4enc (KAESL-OTE1) 비밀번호 힌트.

    func G4DBread(str raw) -> (G4data content)[str name]
    # gen4kdb 텍스트를 읽고 이름-데이터 쌍의 해시맵을 반환.
    func G4DBwrite((G4data content)[str name] data) -> str
    # 이름-데이터 쌍의 해시맵에서 gen4kdb 텍스트 생성.
    struct G4data
        func Set(interface data) -> error
        # 바이트열/문자열/정수/실수 값으로 데이터 설정.

        .ByteV byte[] # 바이트 값 저장 위치
        .StrV str # 문자열 값 저장 위치
        .IntV int # 정수 값 저장 위치
        .FloatV float # 실수 값 저장 위치
        .Dtype rune # 데이터 값 종류. ('b', 's', 'i', 'f', 'n')

    struct G4io
        func OpenB(byte[] raw, bool isreader)
        # raw를 내부 버퍼로 하는 B모드 읽기/쓰기 구조체 설정.
        func OpenF(str path, bool isreader) -> error
        # path 경로로 F모드 읽기/쓰기 구조체 설정.
        func CloseB() -> byte[]
        # B모드 읽기/쓰기를 종료하고 내부 버퍼 반환.
        func CloseF()
        # F모드 내부 파일 포인터를 닫음.
        func Size() -> int
        # 버퍼/파일의 크기 반환.
        func Seek(int pos)
        # 버퍼/파일의 읽기 기준위치 설정.
        func Read(int size) -> byte[]
        # size 크기만큼 읽고 바이트열 반환.
        # 만약 읽기 데이터가 부족하다면 size보다 작은 크기로 반환.
        # 읽기 기준위치는 같은 크기만큼 뒤로 이동.
        func Write(byte[] chunk)
        # 바이트열을 버퍼에 이어붙임/파일에 쓰기.

        .IsBin bool # 바이너리 모드 여부. True면 binary, False면 file 모드.
        .IsReader # 읽기 용도인지 여부. True면 읽기만 가능, False면 쓰기만 가능.

    struct G4kaesall
        func EnBin(byte[] pw, byte[] kf, byte[] data) -> (byte[], error)
        # 바이트열 암호화. 암호화 결과물 반환. (사진 위장헤더 포함)
        func EnFile(byte[] pw, byte[] kf, str path) -> (str, error)
        # 파일 암호화. 암호화 대상 파일과 같은 위치에 새 파일 생성. 암호파일 경로 반환.
        func DeBin(byte[] pw, byte[] kf, byte[] data) -> (byte[], error)
        # 바이트열 복호화. 복호화 결과물 반환. (원본이름 복구기능 없음)
        func DeFile(byte[] pw, byte[] kf, str path) -> (str, error)
        # 파일 복호화. 복호화 대상 파일과 같은 위치에 새 파일 생성. 원본파일 경로 반환.
        func ViewBin(byte[] data) -> error
        # 암호화 바이트열의 정보 해석.
        func ViewFile(str path) -> error
        # 암호화 파일의 정보 해석.

        .Hint byte[] # 비밀번호 힌트.

    struct G4kaesfunc
        func Encrypt(byte[] mkey) -> error
        # 48B mkey로 암호화.
        func Decrypt(byte[] mkey) -> error
        # 48B mkey로 복호화.

        .Inbuf G4io # 입력버퍼 (읽기전용), Open/Close 작업은 모듈 사용자가 처리해야 함.
        .Exbuf G4io # 출력버퍼 (쓰기전용), Open/Close 작업은 모듈 사용자가 처리해야 함.

    func InitKV4(str path) -> g4kv4*
    # 클러스터 경로로 path를 설정하고 내부 필드를 초기화한 구조체 반환.
    struct g4kv4
        func View() -> error
        # 암호화된 클러스터 정보 해석.
        func Read(byte[] pw, byte[] kf, str newpath) -> error
        # 암호화된 클러스터를 newpath 폴더 아래에 복호화. bin/main 폴더가 생성됨.
        func Write(byte[] pw, byte[] kf, str tgtpath) -> error
        # 일반 폴더 tgtpath에서 암호화 클러스터 생성. 클러스터 main 안에 tgtpath가 있는 구조로 생성됨.

        .Path str # 클러스터 경로. (클러스터 읽기/쓰기 모두)
        .Hint byte[] # 클러스터 비밀번호 힌트.

이 라이브러리는 오래된 데이터 형식의 읽고쓰기를 지원하기 위해 만들어졌으며,
시험 버전이기에 충분한 테스트를 거치지 않았습니다.
사용 시 원본 python 코드를 기준으로 하십시오.
(원본 py코드가 잘못 짜인 오류까지도 따라서 구현함)

!!! 폴더 자동초기화 주의. G3kzip.Unpack: "./temp261/", G3picre: "./temp365/" !!!

!! G2enc 패딩, G3kaes 보조 키 생성 과정의 원본 python 코드가 잘못 구현되어 있음. !!

G1ENC {
    magicnum : 4B; .kos // 버전 식별자.
    salt : 40B // salt 바이트.
    pwhash : 32B // 비밀번호 해시.
    hint : 324B // 최대 324B 길이 비밀번호 힌트.

    data : nB // 암호화 데이터의 길이는 원본과 같음.
}

G2ENC {
    magicnum : 4B; kos2 // 버전 식별자.
    salt : 80B // salt 바이트.
    pwhash : 64B // 비밀번호 해시.
    hint : 600B // 최대 600B 길이 비밀번호 힌트.
    encname : 256B // 암호화된 파일명. 0x00으로 패딩됨.
    namelen : 2B // 원본파일명 길이.
    namemode : 2B // 이름숨김 모드에 따라 hi 또는 op.
    header md5 : 16B // 앞 7개 항목의 md5 해시값.

    data : 16nB // 암호화 데이터는 16배수로 패딩됨.
    // 파이썬 버전은 원본길이가 16배수인 경우 0x16이 아닌 0x00을 16개 패딩하는 오류가 존재.
}

G3KDB {
    유니코드 기반 구조적 데이터 형식.
    [이름] {값}. 중괄호 안에 다른 구조체가 포함될 수 있음.
    값으로는 정수, 실수, 문자열, 리스트가 가능.
    이름에 # 포함 시 주석. 문자열 포메팅은 ## : #, #" : ", #s : 공백, #n : 줄바꿈
    [x0]{[x1]{123}[x2]{45.6000}[x3]{"가나다abc"}[y0]{[y1]{128,"###"#s#n"}}}
}

G3KZIP {
    prehead : 1024nB // 위장용 사진 헤더.
    mainhead : 18B {
        magicnum : 4B; KTS2 // 버전 식별자.
        reserved : 2B // 예약됨. (현재 사용하지 않음)
        chunknum : 3B // 청크 개수.
        typelen : 1B // type 길이.
        sizelen : 1B // size 길이.
        namelen : 3B // name 길이.
        crc32data : 4B // 모든 subheader을 이어붙인 것의 CRC32 값.
    }

    chunk *N {
        subheader : (type + size + name)B // S (폴더 구조 정보), F (파일 바이너리).
        data : nB // F 모드는 파일 데이터.
        // S 모드는 내부 폴더를 리스트로 나열함. [folders]{[data]{"x0","x0/x1"}}
    }
    trash : nB // 모든 청크 종료 후에 오는 값은 쓰레기값.
}

G3KAES {
    prehead : 1024nB // 위장용 사진 헤더.
    metadata : 18B {
        magicnum : 4B; KES3 // 버전 식별자.
        reserved : 2B // 예약됨. (현재 사용하지 않음)
        mhsize : 4B // mainhead 크기.
        shsize : 4B // subhead 크기.
        crc32data : 4B // mainhead + subhead의 CRC32 값.
    }
    mainhead : nB // 병렬처리와 파일이름에 대한 데이터가 존재.
    // all-mode : (core, chunksize, ckeydt, salt, pwhash, hint, tkeydt, enctitle)
    // func-mode : (core, chunksize, iv, ckeydt)
    subhead : nB // 현재 사용하지 않는 보조헤더.

    data : 16nB // 데이터는 chunksize 크기로 분할되며, 한번에 core개 만큼 병렬처리됨.
    // chunksize는 16의 배수로 설정되기에, padding은 마지막 청크에 한해 계산됨.
}

G3KPIC {
    사진 바이너리에 데이터를 분할하여 숨김.
    1x 모드는 원본 사진과 무관하게 사진 1 바이트당 데이터 1 바이트 할당.
    2x 모드는 원본 사진의 하위 4비트를 이용해 사진 2 바이트당 데이터 1 바이트 할당.
    RGB 256 값 중 16으로 나눈 몫은 유지하고 나머지를 이용해 인코딩.
    데이터를 16으로 나눈 몫과 나머지를 순서대로 (빅엔디안) 사진과 결합.

    사진 파일은 가로세로가 모두 4의 배수여야 하며, 색이 24비트여야 함.
    이때 2x 모드는 사진 한 장당 row * col * 3 / 2 만큼 저장 가능.
    데이터가 사진 저장 크기와 맞지 않는 경우에는 뒤에 0 패딩을 하기에 kzip같은 방식과 조합하여야 함.
}

G3PICRE {
    앞에 일반 사진 파일이 오게 하고, 뒤에 zip 파일을 위치시킴.
    zip 헤더를 적절히 조절하여 사진과 압축파일 모두로 작동할 수 있게 함.

    zip {
        chunk *N {
            local file head
            compressed file data
        }
        central head *N // 각 local file head의 시작 오프셋이 존재.
        zip mainhead *1 // central head의 시작 오프셋이 존재.
    }
    앞에 사진파일을 추가하며 zip mainhead, central head의 오프셋을 수정하는 원리.
}

G4ENC {
    magicnum : 4B; OTE1 // 버전 식별자.
    hintlen : 2B // 힌트 길이.
    hint : nB // 힌트 바이트열.
    salt : 32B // 비밀번호 salt.
    pwhash : 32B // 비밀번호 hash.
    cketdata : 128B // encrypted content key.
    iv : 16B // plain iv 바이트열.

    data : 16nB // G3KAES와 비슷하게 병렬처리됨. (core, chunk)는 (32, 128k)로 고정.
}

G4KDB {
    바이트열, 문자열, 정수, 실수의 간단한 텍스트 형식 기록.
    모두 대문자로 기록되며, 식별자가 모두 대문자. 각 항목 사이 구분은 줄바꿈 문자.
    바이트열/문자열은 hex 출력값을, 정수/실수는 그대로 문자열화한 값을 적음.
    형식 : name(type)data\n
    DATA0(BYTES)16C3\nDATA1(STR)414243\nDATA2(INT)42\nDATA3(FLOAT)0.123
}

G4KAES {
    prehead : 128nB // 위장용 사진 헤더.
    magicnum : 5B; KAES4 // 버전 식별자.
    mhsize : 3B // mainhead 길이.
    mainhead : nB // bin/file 여부, all/func 모드에 따라 포함된 필드가 다름.
    // all-mode-bin : (MODE, SALT, PWH, CKDT, HINT), all-mode-file : (TKDT, NMDT) 추가됨.
    // func-mode : (MODE, CKDT)

    data : 16nB // G3KAES와 비슷하게 병렬처리됨. (core, chunk)는 (32, 128k)로 고정.
}

G4KV4 {
    header {
        magicnum : 4B; KV4H // 버전 식별자.
        mhsize : 8B // mainhead 크기.
        mainhead : nB // 6개 필드를 가진 mainhead.
        fssize : 8B // filesys 크기.
        filesys : 16nB // filesys 크기.
        fksize : 8B // filekey 크기.
        filekey : 16nB // filekey 크기.
    }

    mainhead {
        MODE, SALT, PWH, AKDT, TKDT, HINT 필드를 가짐.
        akey로 filekey가 암호화되고, tkey로 filesys가 암호화됨.
    }

    filesys {
        가상 파일 시스템은 휴지통 역할의 bin과 저장소 역할의 main으로 나눠 구현됨.
        각 폴더는 하위 폴더와 파일을 가질 수 있으며, 더 하위일수록 "깊이"가 깊어짐.
        이 관계는 파일과 폴더 각 항목을 줄바꿈 문자로 구분하는 이진 파일 형태로 저장됨.
        형식 : 깊이 + 식별자 + 이름 + 추가데이터
        식별자는 폴더는 #, 일반파일은 $, G3KZIP 처리된 폴더파일은 &.
        폴더는 추가데이터가 없으며, 파일형의 경우 슬래시 + fptr이 붙음.

        0#bin
        1#trashdir
        1&trashfile/............
        0#main
        1#upfolder
        2#midfolder
        3#downfolder
        3$data/............
        2$data/............
    }

    filekey {
        모든 파일은 고유한 fptr : 12B과 fkey : 48B를 가짐.
        fptr을 오름차순 정렬되는 순서로 fptr + fkey : 60B씩 이어붙여
        filekey 바이너리를 생성함. (즉 filekey의 길이는 60 * 파일 수)
    }

    클러스터란 헤더, 청크(암호화된 파일들이 모여있는 폴더)가 모여있는 폴더를 말함.
    각 청크 안의 파일 개수는 maxnum에 의해 결정됨.
    모든 클러스터 안의 파일은 대응되는 fptr을 가지며,
    폴더는 filesys에 적힌 가상의 관계도에 의해 생성됨.
}
