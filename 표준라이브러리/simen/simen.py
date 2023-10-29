# test535 : stdlib5 khash/simen

class toolbox:

    def __init__(self):
        prebox = [14, 4, 13, 1, 2, 15, 10, 6, 8, 3, 11, 9, 5, 12, 7, 0]
        self.sbox = [ 16 * prebox[x // 16] + prebox[x % 16] for x in range(0, 256) ] # sbox 0~255 -> r
        boolize = lambda x: [ False if x[t - 8] == '0' else True for t in range(0, 8) ]
        self.bit = [ boolize( '00000000' + bin(x)[2:] ) for x in range(0, 256) ] # 0 ~ 255 -> TF
        self.invsbox = [0] * 256 # sbox inv r -> 0~255
        for i in range(0, 256):
            self.invsbox[ self.sbox[i] ] = i

    def shift(self, array, value): # shift 밀기 연산, 0 ~ 31
        add = value % 4
        ptr = value // 4
        if ptr < 4:
            temp = [ array[4 * ptr], array[4 * ptr + 1], array[4 * ptr + 2], array[4 * ptr + 3] ] * 2
            a = temp[add]
            b = temp[add + 1]
            c = temp[add + 2]
            d = temp[add + 3]
            array[4 * ptr] = a
            array[4 * ptr + 1] = b
            array[4 * ptr + 2] = c
            array[4 * ptr + 3] = d
            
        else:
            ptr = ptr - 4
            temp = [ array[ptr], array[ptr + 4], array[ptr + 8], array[ptr + 12] ] * 2
            a = temp[add]
            b = temp[add + 1]
            c = temp[add + 2]
            d = temp[add + 3]
            array[ptr] = a
            array[ptr + 4] = b
            array[ptr + 8] = c
            array[ptr + 12] = d

    def invshift(self, array, value): # shift inv 연산, 0 ~ 31
        add = value % 4
        ptr = value // 4
        self.shift(array, 4 * ptr + (4 - add) % 4)

    def hash(self, data):
        # 8바이트 패딩
        length = len(data)
        temp = (8 - length % 8) % 8
        data = data + bytes( [temp] ) * temp
        data = [ data[8 * x:8 * x + 8] for x in range(0, len(data) // 8) ]

        iv = [66, 75, 235, 179, 145, 234, 180, 128] # 8B iv
        logic0 = lambda x, y: [ ( not( x[t] ) and y[t] ) and ( x[t] or y[t] ) for t in range(0, 8) ]
        logic1 = lambda x, y: [ ( x[t] and not( y[t] ) ) or ( x[t] ^ y[t] ) for t in range(0, 8) ]
        logic2 = lambda x, y: [ ( not( x[t] ) or y[t] ) or ( x[t] and y[t] ) for t in range(0, 8) ]
        logic3 = lambda x, y: [ ( x[t] ^ y[t] ) and ( x[t] or not( y[t] ) ) for t in range(0, 8) ]
        numerize = lambda x: sum( [ 2 ** (7 - t) if( x[t] ) else 0 for t in range(0, 8) ] )

        for i in data: # data chunk 반복
            iv = [ self.bit[ iv[x] ^ i[x] ] for x in range(0, 8) ] # iv ^ chunk
            matrix = [0] * 16
            for j in range(0, 4): # logic 연산
                matrix[4 * j] = numerize( logic0( iv[j], iv[7 - j] ) )
                matrix[4 * j + 1] = numerize( logic1( iv[j], iv[7 - j] ) )
                matrix[4 * j + 2] = numerize( logic2( iv[j], iv[7 - j] ) )
                matrix[4 * j + 3] = numerize( logic3( iv[j], iv[7 - j] ) )
            for j in i: # shift
                self.shift(matrix, j % 32)
                self.shift(matrix, j // 8)
            matrix = [self.bit[ self.sbox[x] ] for x in matrix] # sbox 치환 후 matrix 합치기
            iv[0] = numerize( [ matrix[14][t] and matrix[4][t] for t in range(0, 8) ] )
            iv[1] = numerize( [ matrix[13][t] ^ matrix[1][t] for t in range(0, 8) ] )
            iv[2] = numerize( [ matrix[2][t] ^ matrix[15][t] for t in range(0, 8) ] )
            iv[3] = numerize( [ matrix[10][t] and matrix[6][t] for t in range(0, 8) ] )
            iv[4] = numerize( [ matrix[8][t] ^ matrix[3][t] for t in range(0, 8) ] )
            iv[5] = numerize( [ matrix[11][t] or matrix[9][t] for t in range(0, 8) ] )
            iv[6] = numerize( [ matrix[5][t] or matrix[12][t] for t in range(0, 8) ] )
            iv[7] = numerize( [ matrix[7][t] ^ matrix[0][t] for t in range(0, 8) ] )

        return bytes(iv)

    def setkey(self, data):
        # 8바이트 패딩
        length = len(data)
        if length == 0:
            raise Exception("keylength0error")
        temp = (8 - length % 8) % 8
        data = data + bytes( [temp] ) * temp
        data = [ data[8 * x:8 * x + 8] for x in range(0, len(data) // 8) ]

        self.key = [0] * len(data)
        reverse = lambda x: list( self.hash( bytes( reversed(x) ) ) )
        for i in range( 0, len(data) ):
            k0a = [ x for x in data[i] ]
            k0b = reverse(k0a)
            k0b = k0b[5:8] + k0b[0:5]
            k1a = [ k0a[x] ^ k0b[x] for x in range(0, 8) ]
            k1b = reverse(k1a)
            k1b = k1b[3:8] + k1b[0:3]
            k2a = [ k1a[x] ^ k1b[x] for x in range(0, 8) ]
            k2b = reverse(k2a)
            k2b = k2b[7:8] + k2b[0:7]
            k3a = [ k2a[x] ^ k2b[x] for x in range(0, 8) ]
            k3b = reverse(k3a)
            k3b = k3b[1:8] + k3b[0:1]
            k4a = [ k3a[x] ^ k3b[x] for x in range(0, 8) ]
            k4b = reverse(k4a)
            k4b = k4b[5:8] + k4b[0:5]
            self.key[i] = [k1a, k1b, k2a, k2b, k3a, k3b, k4a, k4b]

    def encrypt(self, data):
        length = len(data)
        temp = 16 - length % 16
        data = data + bytes( [temp] ) * temp
        temp = len(data) // 16
        data = [ data[16 * x:16 * x + 16] for x in range(0, temp) ]
        output = [0] * temp

        for i in range(0, temp):
            matrix = list( data[i] ) # 4x4
            key = self.key[ i % len(self.key) ] # round key [ []*8 ]

            for j in range(0, 4): # 4 round
                
                for k in range(0, 16): # xor
                    if k % 4 < 2:
                        matrix[k] = matrix[k] ^ key[2 * j][k // 2 + k % 2]
                    else:
                        matrix[k] = matrix[k] ^ key[2 * j + 1][k // 2 + k % 2 - 1]

                for k in range(0, 8): # 8회 shift
                    self.shift(matrix, ( key[2 * j][k] + key[2 * j + 1][k] ) % 32)

                matrix = [self.sbox[x] for x in matrix] # sbox 치환

            output[i] = bytes(matrix)

        return b''.join(output)

    def decrypt(self, data):
        length = len(data)
        temp = length // 16
        data = [ data[16 * x:16 * x + 16] for x in range(0, temp) ]
        output = [0] * temp

        for i in range(0, temp):
            matrix = list( data[i] ) # 4x4
            key = self.key[ i % len(self.key) ] # round key [ []*8 ]

            for j in range(0, 4): # 4 round
                j = 3 - j

                matrix = [self.invsbox[x] for x in matrix] # sbox 치환

                for k in range(0, 8): # 8회 shift
                    k = 7 - k
                    self.invshift(matrix, ( key[2 * j][k] + key[2 * j + 1][k] ) % 32)

                for k in range(0, 16): # xor
                    if k % 4 < 2:
                        matrix[k] = matrix[k] ^ key[2 * j][k // 2 + k % 2]
                    else:
                        matrix[k] = matrix[k] ^ key[2 * j + 1][k // 2 + k % 2 - 1]
            output[i] = bytes(matrix)

        temp = output[-1]
        output[-1] = temp[ 0:-temp[-1] ]
        return b''.join(output)
