# test 631 : stdlib5.kdb

class toolbox: # for lite setting files
    def __init__(self):
        self.name = dict() # (int index)[str name]
        self.type = [ ] # type int[]
        self.ptr = [ ] # pointer int[]
        self.fmem = [ ] # mem float[]
        self.cmem = [ ] # mem complex[]
        self.bmem = [ ] # mem bytes[]

        self.working = [ ] # current names list

    # parse string
    def read(self, raw):
        # pass id setter nonstr str sharp end
        self.working = [ ]
        order = [ x + '\n' for x in raw.split('\n') ]

        for i in order:
            status, mem, name, var = 'pass', [ ], "", ""
            for j in i:
                if status == 'pass':
                    if j not in [' ', '\r']:
                        if j in ['\n', ';']:
                            status, mem, name, var = 'pass', [ ], "", ""
                        elif j == '=':
                            raise Exception(f"invalid key : {i}")
                        else:
                            status = 'id'
                            mem.append(j)
                elif status == 'id':
                    if j in ['\n', ';']:
                        status, mem, name, var = 'pass', [ ], "", ""
                    elif j not in [' ', '\r']:
                        if j == '=':
                            status = 'setter'
                            name = ''.join(mem)
                            mem = [ ]
                        else:
                           mem.append(j)
                elif status == 'setter':
                    if j not in [' ', '\r']:
                        if j in ['\n', ';']:
                            raise Exception(f"invalid value : {i}")
                        elif j == '"':
                            status = 'str'
                            mem.append(j)
                        else:
                            status = 'nonstr'
                            mem.append(j)
                elif status == 'nonstr':
                    if j not in [' ', '\r']:
                        if j == '\n':
                            var = ''.join(mem)
                            self.add(name, var, '\n')
                            status, mem, name, var = 'pass', [ ], "", ""
                        elif j == ';':
                            var = ''.join(mem)
                            self.add(name, var, ';')
                            status, mem, name, var = 'pass', [ ], "", ""
                        else:
                            mem.append(j)
                elif status == 'str':
                    if j == '\n':
                        raise Exception(f"invalid value : {i}")
                    elif j == '#':
                        status = 'sharp'
                    elif j == '"':
                        mem.append(j)
                        var = ''.join(mem)
                        status = 'end'
                    else:
                        mem.append(j)
                elif status == 'sharp':
                    if j == '#':
                        mem.append('#')
                        status = 'str'
                    elif j == 's':
                        mem.append(' ')
                        status = 'str'
                    elif j == 'n':
                        mem.append('\n')
                        status = 'str'
                    elif j == '"':
                        mem.append('"')
                        status = 'str'
                    else:
                        raise Exception(f"invalid escaping : {i}")
                elif status == 'end':
                    if j == '\n':
                        self.add(name, var, '\n')
                        status, mem, name, var = 'pass', [ ], "", ""
                    elif j == ';':
                        self.add(name, var, ';')
                        status, mem, name, var = 'pass', [ ], "", ""

    # return current DB str
    def write(self):
        out = [""] * len(self.type)
        for i in self.name:
            toadd = i
            num = self.name[i]
            tp = self.type[num]
            ptr = self.ptr[num]
            if tp // 16 == 0:
                end = '\n'
            else:
                end = '; '
            tp = tp % 16

            toadd = toadd + ' = '

            if tp == 0:
                toadd = toadd + self.conv(None)
            elif tp == 1:
                if ptr == 0:
                    toadd = toadd + self.conv(False)
                else:
                    toadd = toadd + self.conv(True)
            elif tp == 2:
                toadd = toadd + self.conv(ptr)
            elif tp == 3:
                toadd = toadd + self.conv( self.fmem[ptr] )
            elif tp == 4:
                toadd = toadd + self.conv( self.cmem[ptr] )
            elif tp == 5:
                toadd = toadd + self.conv( self.bmem[ptr] )
            elif tp == 6:
                temp = str(self.bmem[ptr], encoding='utf-8')
                toadd = toadd + self.conv(temp)
            else:
                raise Exception(f"invalid type : {tp}")
            toadd = toadd + end
            out[num] = toadd
        return ''.join(out)

    # get [index, type, ptr, value] by name
    def get(self, name):
        name = name.replace('/', '.')
        num = self.name[name]
        tp = self.type[num] % 16
        ptr = self.ptr[num]

        if tp == 0:
            var = None
            t = 'None'
        elif tp == 1:
            if ptr == 0:
                var = False
            else:
                var = True
            t = 'bool'
        elif tp == 2:
            var = ptr
            t = 'int'
        elif tp == 3:
            var = self.fmem[ptr]
            t = 'float'
        elif tp == 4:
            var = self.cmem[ptr]
            t = 'complex'
        elif tp == 5:
            var = self.bmem[ptr]
            t = 'bytes'
        elif tp == 6:
            var = str(self.bmem[ptr], encoding='utf-8')
            t = 'str'
        else:
            raise Exception(f"invalid type : {tp}")

        return [num, t, ptr, var]

    # revice data by name
    def fix(self, name, data):
        name = name.replace('/', '.')
        num = self.name[name]
        end = self.type[num] // 16
        end = end * 16

        if data == None:
            self.type[num] = end + 0
            self.ptr[num] = 0
        elif type(data) == bool:
            self.type[num] = end + 1
            if data:
                self.ptr[num] = 1
            else:
                self.ptr[num] = 0
        elif type(data) == int:
            self.type[num] = end + 2
            self.ptr[num] = data
        elif type(data) == float:
            self.type[num] = end + 3
            self.ptr[num] = len(self.fmem)
            self.fmem.append(data)
        elif type(data) == complex:
            self.type[num] = end + 4
            self.ptr[num] = len(self.cmem)
            self.cmem.append(data)
        elif type(data) == bytes:
            self.type[num] = end + 5
            self.ptr[num] = len(self.bmem)
            self.bmem.append(data)
        elif type(data) == str:
            self.type[num] = end + 6
            self.ptr[num] = len(self.bmem)
            self.bmem.append( bytes(data, encoding='utf-8') )
        else:
            raise Exception(f"invalid type : {type(data)}")

    # convert data -> kformat str
    def conv(self, data):
        if data == None:
            return 'None'
        elif type(data) == bool:
            if data:
                return 'True'
            else:
                return 'False'
        elif type(data) == int:
            return str(data)
        elif type(data) == float:
            return str(data)
        elif type(data) == complex:
            temp = str(data)
            if temp[0] == '(':
                temp = temp[1:-2] + 'j'
            else:
                temp = temp[0:-1] + 'j'
            return temp
        elif type(data) == bytes:
            temp = [0] * len(data)
            for i in range( 0, len(data) ):
                if data[i] > 15:
                    temp[i] = str( hex( data[i] ) )[2:4]
                else:
                    temp[i] = '0' + str( hex( data[i] ) )[2]
            return "'" + ''.join(temp) + "'"
        elif type(data) == str:
            data = data.replace('#', '##')
            data = data.replace(' ', '#s')
            data = data.replace('\n', '#n')
            data = data.replace('"', '#"')
            return '"' + data + '"'
        else:
            raise Exception(f"invalid type : {type(data)}")

    # add DB by name, var, end
    def add(self, name, var, end):
        name = name.replace('/', '.')
        num = 0
        while name[num] == '.':
            num = num + 1
        self.working = self.working[0:num] + [ name[num:] ]
        name = '.'.join(self.working)
        
        if end == '\n':
            tp = 0
        else:
            tp = 16

        if var == 'None':
            self.name[name] = len(self.type)
            self.type.append(tp + 0)
            self.ptr.append(0)
        elif var == 'True':
            self.name[name] = len(self.type)
            self.type.append(tp + 1)
            self.ptr.append(1)
        elif var == 'False':
            self.name[name] = len(self.type)
            self.type.append(tp + 1)
            self.ptr.append(0)
        elif var[0] == '"':
            self.name[name] = len(self.type)
            self.type.append(tp + 6)
            self.ptr.append( len(self.bmem) )
            self.bmem.append( bytes(var[1:-1], 'utf-8') )
        elif var[0] == "'":
            self.name[name] = len(self.type)
            self.type.append(tp + 5)
            self.ptr.append( len(self.bmem) )
            self.bmem.append( bytes.fromhex( var[1:-1].lower() ) )
        elif var[-1] == 'j':
            self.name[name] = len(self.type)
            self.type.append(tp + 4)
            self.ptr.append( len(self.cmem) )
            self.cmem.append( complex( var.replace('i', 'j') ) )
        elif '.' in var:
            self.name[name] = len(self.type)
            self.type.append(tp + 3)
            self.ptr.append( len(self.fmem) )
            self.fmem.append( float(var) )
        else:
            self.name[name] = len(self.type)
            self.type.append(tp + 2)
            self.ptr.append( int(var) )

    # import and set by input pylist
    def imp(self, arr):
        for i in arr:
            name = i[0]
            data = i[1]
            end = i[2]
            
            num = len(self.type)
            self.name[name] = num
            if end == '\n':
                end = 0
            else:
                end = 16
            self.type.append(0)
            self.ptr.append(0)

            if data == None:
                self.type[num] = end + 0
                self.ptr[num] = 0
            elif type(data) == bool:
                self.type[num] = end + 1
                if data:
                    self.ptr[num] = 1
                else:
                    self.ptr[num] = 0
            elif type(data) == int:
                self.type[num] = end + 2
                self.ptr[num] = data
            elif type(data) == float:
                self.type[num] = end + 3
                self.ptr[num] = len(self.fmem)
                self.fmem.append(data)
            elif type(data) == complex:
                self.type[num] = end + 4
                self.ptr[num] = len(self.cmem)
                self.cmem.append(data)
            elif type(data) == bytes:
                self.type[num] = end + 5
                self.ptr[num] = len(self.bmem)
                self.bmem.append(data)
            elif type(data) == str:
                self.type[num] = end + 6
                self.ptr[num] = len(self.bmem)
                self.bmem.append( bytes(data, encoding='utf-8') )
            else:
                raise Exception(f"invalid type : {type(data)}")

    # export to precise pylist
    def exp(self):
        out = [0] * len(self.type) # [fullname str, data, end]
        for i in self.name:
            temp = ["", "", ""]
            num = self.name[i]
            tp = self.type[num]
            ptr = self.ptr[num]

            temp[0] = i
            if tp // 16 == 0:
                temp[2] = '\n'
            else:
                temp[2] = ';'
            tp = tp % 16

            if tp == 0:
                temp[1] = None
            elif tp == 1:
                if ptr == 0:
                    temp[1] = False
                else:
                    temp[1] = True
            elif tp == 2:
                temp[1] = ptr
            elif tp == 3:
                temp[1] = self.fmem[ptr]
            elif tp == 4:
                temp[1] = self.cmem[ptr]
            elif tp == 5:
                temp[1] = self.bmem[ptr]
            elif tp == 6:
                temp[1] = str(self.bmem[ptr], encoding='utf-8')
            else:
                raise Exception(f"invalid type : {tp}")
            
            out[num] = temp
        return out
