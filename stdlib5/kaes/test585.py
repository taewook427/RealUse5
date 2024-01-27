def main():
    print()
    num = 7
    data = [0.0] * num
    for i in range(0, num):
        data[i] = float(input(f"속도 데이터 {i} 입력 : "))
    print(f"평균값 : {sum(data) / num}")
    print(f"최댓값 : {max(data)}")
    print(f"최솟값 : {min(data)}")
    return input("input y to exit... ")

temp = True
while temp:
    if main() == "y":
        temp = False
