stdlib5.kzip [win & linux] : 파일/폴더의 병렬적 단일 파일 패키저.

<python>
func dozip(str[] tgts, str mode, str path)
# 타겟 리스트의 항목을 패키징해 path 경로에 단일 파일로 생성.
func unzip(str path, str export, bool chkerr)
# path 경로의 파일을 언패킹하여 export 폴더에 내용물 생성.

<go>
func Dozip(str[] tgts, str mode, str path) -> error
# 타겟 배열의 항목을 패키징해 path 경로에 단일 파일로 생성.
func Unzip(str path, str export, bool chkerr) -> error
# path 경로의 파일을 언패킹하여 export 폴더에 내용물 생성.

폴더/파일 상관없이 단일 파일로 패키징 할 수 있습니다.
타겟 리스트 안의 경로는 폴더일 경우 그 폴더와 전체 하위 항목이, 파일인 경우 그 파일이 패키징됩니다.
이후 풀기 경로의 폴더 안에 타겟 리스트의 항목들이 다시 생성됩니다.

!!! 예약된 폴더명 경고 !!!
최상위 폴더 ("D:/", "/" 등)는 기록 시 "LargeVolume{N}/" 꼴의 이름으로 바뀌어 기록됩니다.
해당 이름을 사용하는 폴더는 경로 간섭의 위험이 있습니다.

KSC5 형식의 pre-header을 결정하는 패키징 함수의 mode는 "webp"/"png"/""을 지원하며,
각각 webp/png 형식의 사진과 사진헤더 없음을 의미합니다.
언패킹 함수에서 예외체크가 참이라면, subtype/CRC32/chunk개수 3가지 항목이 잘못되었는지 검사합니다.
패킹/언패킹 과정에서 새로 생성되는 파일/폴더의 경로에 있는 내용은 자동으로 초기화됩니다.
