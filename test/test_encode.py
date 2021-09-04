# -*- coding: utf-8 -*-
import struct
import time

__MAGIC__ = b'\xef\xbe\xad\xde'
__FORMAT_VERSION__ = b'\x02'

__FIELDS_ORIGIN__ = 1 << 0
__FIELDS_ANGLES__ = 1 << 1
__FIELDS_VELOCITY__ = 1 << 2

with open('output/round2/olofmeister.rec', 'rb') as iFile:
    # Step 1: valid check
    _buffer = iFile.read(4)
    assert _buffer == __MAGIC__
    
    # Step 2: version check
    _buffer = iFile.read(1)
    assert _buffer == __FORMAT_VERSION__

    # Step 3: endtime timestamp
    _buffer = iFile.read(4)
    timestamp, = struct.unpack('i', _buffer)
    dt = time.localtime(timestamp)
    print(dt)

    # Step 4: name length
    _buffer = iFile.read(1)
    nameLength, = struct.unpack('b', _buffer)
    print(f'nameLength: {nameLength}')

    # Step 5: record name
    _buffer = iFile.read(nameLength)
    recName =  _buffer.decode()
    print(f'recName: {recName}')
    
    # Step 6: initial position
    initPos = []
    for idx in range(3):
        _buffer = iFile.read(4)
        ipos, = struct.unpack('f', _buffer)
        initPos.append(ipos)
    
    # Step 7: initial angle
    initAngle = []
    for idx in range(2):
        _buffer = iFile.read(4)
        iangle, = struct.unpack('f', _buffer)
        initAngle.append(iangle)
    
    print(f'setpos {initPos[0]} {initPos[1]} {initPos[2]}; setang {initAngle[0]} {initAngle[1]} 0')

    # Step 8: total tick
    _buffer = iFile.read(4)
    totalTick, = struct.unpack('i', _buffer)
    print(f'totalTick: {totalTick}')

    # Step 9: total bookmark(no need)
    _buffer = iFile.read(4)
    totalBM, = struct.unpack('i', _buffer)
    print(f'totalBM: {totalBM}')
    
    # Step 10: read all bookmark
    for bookmark in range(totalBM):
        _buffer = iFile.read(4)  # frame
        _buffer = iFile.read(4)  # additionalTeleportTick
        _buffer = iFile.read(64) # bookmark name, length defined by botmimic: MAX_BOOKMARK_NAME_LENGTH

    # Step 11: read all tick
    for tick in range(totalTick):
        # read 14 items(4bytes each)
        _buffer = iFile.read(4)
        playerButtons, = struct.unpack('i', _buffer)
        _buffer = iFile.read(4)
        playerImpulse, = struct.unpack('i', _buffer)
        actVel, predictVel, predictAng = [], [], []
        for idx in range(3):
            _buffer = iFile.read(4)
            actVel.append(struct.unpack('f', _buffer)[0])
        for idx in range(3):
            _buffer = iFile.read(4)
            predictVel.append(struct.unpack('f', _buffer)[0])
        for idx in range(2):
            _buffer = iFile.read(4)
            predictAng.append(struct.unpack('f', _buffer)[0])
        _buffer = iFile.read(4)
        newWeapon, = struct.unpack('i', _buffer)
        _buffer = iFile.read(4)
        playerSubtype, = struct.unpack('i', _buffer)
        _buffer = iFile.read(4)
        playerSeed, = struct.unpack('i', _buffer)
        _buffer = iFile.read(4)
        addFields, = struct.unpack('i', _buffer)
        
        #if addFields & (__FIELDS_ORIGIN__ | __FIELDS_ANGLES__ | __FIELDS_VELOCITY__):
        if newWeapon != 0:
            print(f'tick: {tick}')
            print(f'playerButtons: {playerButtons}')
            print(f'playerImpulse: {playerImpulse}')
            print(f'actVel: {actVel}')
            print(f'predictVel: {predictVel}')
            print(f'predictAng: {predictAng}')
            print(f'newWeapon: {newWeapon}')
            print(f'playerSubtype: {playerSubtype}')
            print(f'playerSeed: {playerSeed}')
            print(f'addFields: {addFields}')
            print()

        if (addFields &  __FIELDS_ORIGIN__):
            nowOrigin = []
            for idx in range(3):
                _buffer = iFile.read(4)
                nowOrigin.append(struct.unpack('f', _buffer)[0])
        if (addFields &  __FIELDS_ANGLES__):
            nowAngle = []
            for idx in range(3):
                _buffer = iFile.read(4)
                nowAngle.append(struct.unpack('f', _buffer)[0])
        if (addFields &  __FIELDS_VELOCITY__):
            nowVelocity = []
            for idx in range(3):
                _buffer = iFile.read(4)
                nowVelocity.append(struct.unpack('f', _buffer)[0])
