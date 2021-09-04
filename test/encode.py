# -*- coding: utf-8 -*-

import struct
import io
import time
import ujson

__PREFIX__ = b'\xef\xbe\xad\xde'
__FORMAT_VERSION__ = b'\x02'

__FIELDS_ORIGIN__ = 1 << 0
__FIELDS_ANGLES__ = 1 << 1
__FIELDS_VELOCITY__ = 1 << 2

PLAYER = 'ZywOo'

def __encode(jsonFile: io.BufferedReader, outFile: io.BufferedWriter):
    data_dict = ujson.load(jsonFile)

    # st.1 valid check
    outFile.write(__PREFIX__)
    # st.2 format version check
    outFile.write(__FORMAT_VERSION__)
    # st.3 timestamp
    timenow = struct.pack('i', int(time.time()))
    outFile.write(timenow)
    # st.4 name length
    namelen = struct.pack('b', len(PLAYER))
    outFile.write(namelen)
    # st.5 name
    name = bytes(PLAYER, encoding='utf-8')
    outFile.write(name)
    # st.6 initial position and angle
    for _player in data_dict['Frames'][0]['CT']['Players']:
        if _player['Name'] == PLAYER and _player['IsAlive']:
            initPosAng = [
                struct.pack('f', _player['X']),
                struct.pack('f', _player['Y']),
                struct.pack('f', _player['Z']),
                struct.pack('f', _player['ViewX']),
                struct.pack('f', _player['ViewY'])
            ]
            for bPos in initPosAng:
                outFile.write(bPos)
            break
    # st.7 total tick
    tickCount = 0
    for frame in data_dict['Frames']:
        end = True
        for _player in frame['CT']['Players']:
            if _player['Name'] == PLAYER and _player['IsAlive']:
                tickCount += 1
                end = False
                break
        if end:
            break
    btickTotal = struct.pack('i', tickCount)
    outFile.write(btickTotal)
    # st.8 bookmark(no need)
    totalBM = struct.pack('i', 0)
    outFile.write(totalBM)
    # st.9 write all bookmarks(pass)
    # st.10 write all ticks
    snapshotInterval = 1000
    interval = 0
    for tick in range(tickCount):
        for _player in data_dict['Frames'][tick]['CT']['Players']:
            if _player['Name'] == PLAYER and _player['IsAlive']:
                # write 14 basic items
                playerButtons = struct.pack('i', _player['Buttons']) # default
                outFile.write(playerButtons)
                playerImpulse = struct.pack('i', 0) # default
                outFile.write(playerImpulse)
                vel2ang1 = [
                    struct.pack('f', _player['VelocityX']), # actualVelocity
                    struct.pack('f', _player['VelocityY']),
                    struct.pack('f', _player['VelocityZ']),
                    struct.pack('f', _player['VelocityX']), # predictVelocity ???
                    struct.pack('f', _player['VelocityY']),
                    struct.pack('f', _player['VelocityZ']),
                    struct.pack('f', _player['ViewY']), # predictAngle ???
                    struct.pack('f', _player['ViewX'])
                ]
                for item in vel2ang1:
                    outFile.write(item)
                break
        newWeapon = struct.pack('i', 0) # default
        outFile.write(newWeapon)
        playerSubtype = struct.pack('i', 0) # default
        outFile.write(playerSubtype)
        playerSeed = struct.pack('i', 0) # default
        outFile.write(playerSeed)
        # additionalFields = struct.pack('i', __FIELDS_ANGLES__ | __FIELDS_ORIGIN__ | __FIELDS_VELOCITY__) # default
        if interval >= snapshotInterval:
            interval = 0
            additionalFields = struct.pack('i', __FIELDS_ORIGIN__)
            outFile.write(additionalFields)
            addfields = [
                struct.pack('f', _player['X']),
                struct.pack('f', _player['Y']),
                struct.pack('f', _player['Z']),
            ]
            for adds in addfields:
                outFile.write(adds)
        else:
            interval += 1
            additionalFields = struct.pack('i', 0) # default
            outFile.write(additionalFields)

def test_encode(jsonpath: str, outpath: str):
    with open(outpath, 'wb') as oFile:
        with open(jsonpath, 'r') as iFile:
            __encode(iFile, oFile)


def main():
    test_encode('json/demo.json', 'rec/demo.rec')

if __name__ == '__main__':
    main()