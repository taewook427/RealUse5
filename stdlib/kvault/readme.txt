stdlib5.kvault [win & linux] : 파일볼트 시스템.

<go>
func PEVFS_New(str remote, str cluster, int chunksize) -> error
# 빈 remote 폴더에 새 클러스터를 생성함.
func PEVFS_Boot(str desktop, str local, str remote, str blockApath) -> (PEVFS*, byte[], error)
# 이미 존재하는 KV5 클러스터 정보를 읽어와 PEVFS 구조체와 비밀번호 힌트 반환.
func PEVFS_Exit(PEVFS* obj)
# local 폴더를 삭제하고 PEVFS 구조체를 내부 데이터를 지워 정지 상태로 만듦.
func PEVFS_Rebuild(str remote, byte[] pw, byte[] kf) -> error
# 보안강화를 위해 fphy를 초기화하고 블록 C를 모두 재작성함.

struct PEVFS
    func Abort(bool reset, bool abort, bool working) -> (bool, bool)
    # reset이 True면 내부 플래그를 재설정함. 현재 상태 플래그(abort/working)를 반환함.
    func Debug() -> (int[], byte[][])
    # 내부 상태를 디버그용으로 반환. int배열 : [chunksize, blocknum, first fptr].
    # 바이트배열 : [wrsign(8B), salt(64B), pwhash(192B), fsyskey(48B), fkeykey(48B), fphykey(48B)].
    func Log(bool reset) -> str
    # '\n'으로 결합된 로그 반환. reset이 True면 로그를 초기화함.

    func Login(byte[] pw, byte[] kf, int sleeptime) -> error
    # Boot된 클러스터에 PWKF로 로그인하여 정보를 얻어옴. 입출력 오류 시 대기시간(초)를 받음.
    func AccReset(byte[] pw, byte[] kf, byte[] hint) -> error
    # 계정 PWKF 변경. R 모드라도 븦록 A는 재작성됨.
    func AccExtend(byte[] pw, byte[] kf, byte[] hint, str account, bool wrlocked) -> (str, error)
    # Curdir을 새 Rootdir로 한 계정 생성, 새 블럭 A는 desktop에 생성되며 그 경로가 반환됨.
    # wrlocked가 True여야 하위 폴더 중 잠금상태인 폴더도 기록됨.

    func Search(str name) -> str[]
    # Curdir 밑의 폴더/파일 중 이름 기반으로 검색. 와일드카드 문자나 이스케이프 패턴 지원.
    # * : 0+ str, ? : len1 str, %d : int, %s 1+ ascii, %c len1 ascii, %* %? %% : literal.
    func Print(bool wrlocked) -> str
    # 현재 폴더와 하위 항목을 디버그용 문자열로 반환. wrlocked가 True여야 잠금상태인 폴더도 기록됨.
    func Teleport(str path) -> bool
    # Rootdir 기준 path(~/) 폴더로 Curdir 변경 후 True 반환. 만약 폴더 없다면 변경하지 않고 False 반환.
    func NavName() -> (str[], bool[])
    # Curdir의 하위 폴더/파일에 대해 이름과 잠김여부 반환. 파일은 잠김 False로 반환.
    func NavInfo(bool wrlocked) -> (str[], int[], int[], int[], int[])
    # 재귀접근을 통해 더 많은 정보를 반환하는 NavName. ( names, times, size, fptr/islocked(0T 1F), [lowerdir, lowerfile] ).
    # 0번째 항목은 Curdir 자신에 대한 정보, wrlocked가 True여야 잠긴폴더도 정보 반환.

    func ImBin(str name, byte[] data) -> error
    # 이진데이터 가져오기 (curpath/name). 이름이 이미 존재한다면 덮어쓰기됨.
    func ImFiles(str[] paths) -> error
    # 파일들 가져오기, 전체경로를 받아 자동으로 파일명을 잡아줌. 이름이 이미 존재한다면 덮어쓰기됨.
    func ImDir(str path) -> error
    # 폴더 가져오기. 전체경로를 받아 자동으로 폴더명을 잡아줌. 이름이 이미 존재한다면 오류발생.
    func ExBin(str name) -> (byte[], error)
    # 이진데이터 내보내기 (curpath/name).
    func ExFiles(str[] names) -> error
    # 파일들 내보내기. 이름을 받아 해당파일을 바탕화면에 생성 (desktop/kv5export/name).
    func ExDir(str path) -> error
    # 폴더 내보내기. 전체경로를 받아 해당폴더를 바탕화면에 생성 (desktop/kv5export/name/).

    func Delete(str[] paths) -> error
    # 전체경로를 받아 파일/포함된 파일을 모두 삭제함.
    func Move(str[] names, str dst) -> error
    # Curdir 하위 폴더/파일을 dst로 지정된 전체경로 하위로 이동시킴. 특수폴더나 위계문제가 있는 폴더는 이동 불가.
    func Rename(str[] before, str[] after) -> error
    # Curdir 하위 폴더/파일의 이름을 before에서 after로 바꿈. 두 이름은 서로 형식이 일치해야 함.
    func DirNew(str[] names) -> error
    # Curdir 하위에 새 폴더들을 생성, 이름은 길이 1 이상으로 생성조건에 맞아야 함.
    func DirLock(str path, bool islocked, bool sub) -> error
    # 전체경로를 받아 해당폴더의 잠금상태를 islocked로 설정, sub가 True면 하위폴더도 마찬가지로 설정.

    func CluCheck() -> (int, error)
    # Curdir 하위의 모든 파일에 대해, fphy 정상인지 체크하고 비정상 파일개수 반환.
    func CluRestore(bool rename, bool rewrite, bool rebuild) -> (int, error)
    # 클러스터의 모든 파일에 대해 복구/재작성 작업 수행후 비정상 파일개수 반환. rename : 위험/비정상 이름 수정.
    # rewrite : fphy rewriteGC 수행 (미삭제 청크 재기록). rebuild : fkey/fphy 정상 체크하고 plainGC 수행 (연결끊긴 청크 해제).

    .Curpath str # 현재 작업 폴더 경로 (~/).
    .Curdir vdir* # 현재 작업 포인터 (사용할 일 없음).
    .Rootpath str # 루트 경로, RW 모드면 '/', R 모드면 '~/'.
    .Cluster str* # 클러스터 이름.
    .Account str* # 계정 이름, RW 모드면 'root', R 모드면 '~'.

struct Shell
    func Command(str order, str[] option) -> error
    # 명령에 따라 작업을 수행, 비동기 처리와 패닉 복구를 자동으로 지원, IObuf를 통해 입출력.

    .InSys PEVFS* # 내부 PEVFS 구조체.
    .AsyncErr str # 비동기 처리 에러결과물, 문제없다면 빈 문자열로 설정됨.
    .IOstr str[] # 입출력 문자열 버퍼.
    .IObyte byte[][] # 입출력 바이트열 버퍼.

    .FlagWk bool # 비동기 작업 처리 플래그 (읽기전용).
    .FlagRo bool # 읽기전용 클러스터 플래그 (읽기전용).
    .FlagSz bool # 크기/시간 정보 가져오기 플래그.

    .CurPath str # 현재 작업경로 (~/).
    .CurNum int[2] # 직접하위 폴더수/파일수.
    .CurName str[] # 직접하위 폴더/파일 이름.
    .CurLock bool[] # 직접하위 폴더/파일의 잠금 여부.
    .CurTime str[] # 직접하위 폴더/파일의 포매팅된 시간.
    .CurSize int[] # 직접하위 폴더/파일의 크기.

func Test_Basic(str remote) -> (float, float, float)
# 클러스터 생성 / remote 위치에 4GiB 데이터 읽기 / 쓰기 속도 측정.
func Test_IO(str remote) -> (float, float, float)
# 클러스터 로그인 / 4GiB 데이터 내보내기 / 가져오기 속도 측정.
func Test_Multi(str remote) -> (float, float, float)
# 실용사용한도 (폴더 10만, 파일 1000만) 생성 / 읽기 / 쓰기 속도 측정.

! shell 명령어 !
순서 : order / option / IOstr / IObyte / Async
옵션 : T "true" / F "false" / N index of CurName / ... (*n) / -> (ret)

init / flagsz[TF]
new / remote, cluster, csize["small", "standard", "large", "*"]
boot / desktop, local, remote, blockApath / cluster[->], account[->] / hint[->]
exit / / / / FG
rebuild / remote / / pw, kf

abort / reset[TF], abort[TF], working[TF] / / / FG
debug / countLock[TF] / info0[->], info1[->], info2[->], info3[->] / / FG
log / reset[TF] / logdt[->] / / FG
login / sleep["1", "10", "30", "60", "*"] / / pw, kf
reset / / / pw, kf, hint
extend / account, wrlock[TF] / newpath[->] / pw, kf, hint

search / name / result[->] / / FG
print / wrlock[TF] / result[->] / / FG
navigate / path / subdir[->] / / FG
update / / / / FG

imbin / name / / data / BG
imfile / fullpath[...] / / / BG
imdir / fullpath / / / BG
exbin / name[...] / / data[->][...] / BG
exfile / name[...] / / / BG
exdir / fullpath / / / BG

delete / nameN[N][...] / / / BG
move / tgtdir, nameN[N][...] / / / BG
rename / nameN[N][...] / newname[...] / / BG
dirnew / name[...] / / / BG
dirlock / lock[TF], nameN[N][...] / / / BG

check / / result[->] / / BG
restore / mode["rename", "rewrite", "rebuild"] / result[->] / / BG
*

!! KV5 용어설명!!
경로 구분자 : 폴더/파일 단계를 구분하는 글자로, 리눅스와 같이 '/'사용.
이름 : 각 폴더/파일에 붙은 이름으로 '/'과 '\n'을 제외한 유니코드 문자열 사용 가능.
잠금상태 : 폴더에 부여될 수 있는 상태로, 잠긴 폴더는 잠김옵션 False인 기록시도로 기록될 수 없음.
파일 : 실제 바이너리 내용을 가지고 암호화되어 원격 저장소에 보관되는 실체.
폴더 : 하위 폴더/파일을 가지며 KV5 시스템 상에서의 논리적 위계계층으로만 존재함. (이름이 '/'로 끝남.)
청크 : 암호화 파일의 일부 데이터를 가진 동일한 크기의 바이트열. 65536개의 청크가 모여 1 블럭이 됨.
블록 : 같은 목적의 데이터를 모아 물리 파티션 상에 하나의 파일로 존재하는 것. (A/B/C/D 블록 존재.)
클러스터 : 하나의 폴더로 표현되는 전체 저장소, 혹은 그 저장소가 위치한 물리적 폴더.
remote : 클러스터가 위치한 원격 저장소 경로, 보통 네트워크/HDD 저속 드라이브 상에 위치.
local : 입출력에 필요한 임시폴더, 보통 SSD 고속 드라이브 상에 위치.
desktop : 사용자 데이터 입출력을 위해 설정된 OS의 바탕화면 위치.
R 모드 : 읽기전용 모드로, 클러스터 데이터를 수정하지 않음.
RW 모드 : 읽고쓰기 모드로, 로그인 시 블록 A를 재작성함.

블록화 등을 통해 보안 취약점을 최소화하고 클러스터 보호도를 높인 파일볼트 시스템입니다.
kvault는 물리적 파티션(ntfs, ext4 등)위에 위치한 파일 블럭(a/b/c/d) 상에서 논리적으로 구성됩니다.
자세한 구조나 주의사항, 사용 방법은 기술문서를 참조하십시오.
