extension.kvdrive [win & linux] : 암호화 파일 시스템 클라이언트.

<Session>
Abort/Recover : 모든 작업을 중단시키거나 클러스터 연결을 복구합니다.
Clear : 검색결과, 로그, View 섹션의 데이터를 지웁니다.
Debuf Info : 디버그 정보(현재 폴더와 세션 키 정보)를 텍스트 섹션에 표시합니다.

<File>
Import : 파일들 혹은 폴더 하나를 가져옵니다.
Export : 현재 선택한 파일들 혹은 폴더 하나를 바탕화면에 내보냅니다.
View : 24MiB이하의 텍스트나 사진, 이진파일을 보거나 수정합니다.

<Manage>
Recycle : 선택된 항목들을 휴지통(/_BIN/)으로 이동합니다.
Rename : 선택된 항목 하나의 이름을 바꿉니다.
Move/Cancel : 선택된 항목들을 다른 폴더로 이동시킵니다.
뷰어가 Move 모드로 바뀌고, 원하는 폴더로 이동한 뒤 좌측 상단 확인 버튼을 눌러 이동시킵니다.

<Control>
Delete : 휴지통 안의 항목을 영구히 삭제합니다.
Deep Lock : 현재 작업 경로의 폴더와 모든 하위 폴더의 잠금 상태를 바꿉니다.
New : 새 파일 또는 폴더를 생성합니다.
Select/Toggle : 전체선택/해제 기능과 폴더 크기 보기/끄기 기능입니다.

<Advanced>
Reset/Export : 계정 비밀번호를 바꾸거나 클러스터의 일부를 열람할 수 있는 새 계정을 생성할 수 있습니다.
계정 생성 시 include Lock을 해야 잠긴 폴더도 볼 수 있습니다.
Restore Name : 모든 항목의 이름을 NTFS 기준에 어긋나지 않게 바꿉니다. 결과가 0이면 정상입니다.
Restore Data : 지워졌다고 표시되었지만 지워지지 않은 청크를 실제로 지웁니다. (쓰레기 수집)
Restore Struct : 클러스터 구조의 무결성을 검사합니다. 결과가 0이면 정상입니다.
Check : 클러스터 무결성을 확인합니다. 결과가 0이면 정상입니다.

<표시줄-상단>
확인 버튼 : 새로고침 혹은 Move 확인 기능.
뒤로가기 버튼 : 상위 폴더로 이동합니다.
표시줄 : 현재 작업 폴더 혹은 View 대상인 파일의 위치를 표시합니다.
위/아래 버튼 : 항목은 한 번에 12개씩 표시됩니다. 표시 위치를 바꿉니다.

<버튼-Lock/Sel>
폴더 항목 하나의 잠금 상태를 바꾸거나 선택 여부를 바꾸는 버튼입니다.
휴지통(/_BIN/)과 버퍼(/_BUF/)는 선택하거나 잠금 상태를 바꿀 수 없습니다.

<섹션-Search>
현재 작업 경로 안에서 이름을 검색하거나 로그를 확인할 수 있습니다.
로그 새로고침 버튼이 활성화되면 작업 상태가 변화한 것입니다.

<섹션-txt/pic/bin>
텍스트/바이너리를 보거나 수정할 수 있습니다.
사진은 여러 파일을 돌아가며 보여줍니다.

<로그인-MkNew>
클러스터 경로로 빈 폴더를 선택하고 로그인하면 새 클러스터가 생성됩니다.
클러스터 타입은 KV4와 KV5가 있으며, 외장하드/USB는 KV5, 클라우드 마운트는 KV4 사용을 권장합니다.

<로그인-Login>
클러스터 폴더 경로와 계정 파일 경로를 설정하고 비밀번호와 키파일을 입력하여 로그인합니다.
올바른 클러스터 타입을 선택해야 합니다.
Uni 옵션이 켜져 있어야 동시에 여러 KVdrive 프로그램을 실행할 수 있습니다.

<로그인-Rebuild>
기존 클러스터의 보안성 향상을 위해 구조를 변화시키는 모드입니다.

Kviewer5와 동일한 UI로 KV4adv, KV5 암호화 파일 시스템을 접근할 수 있습니다.
새 클러스터 생성 시 호스트 파일 시스템에 따라 크기를 다르게 정하는 것이 좋습니다.
(NTFS/ext4 외장하드 : large, FAT32 USB : standard, 클라우드 마운트 : small)
