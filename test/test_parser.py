# -*- coding: utf-8 -*-
import struct
import io

__HEADER__ = 'HL2DEMO' + chr(0)

demofile = './static/faze-vs-vitality-m1-mirage.dem'

def parseHeader(iFile):
    # 01 Header
    _buffer = iFile.read(8)
    header = _buffer.decode()
    assert header == __HEADER__
    # 02 Demo Protocol [little endian]
    _buffer = iFile.read(4)
    demoProtocol, = struct.unpack('<i', _buffer)
    print(f'Header=>DemoProtocol: {demoProtocol}')
    # 03 Network Protocol [little endian]
    _buffer = iFile.read(4)
    networkProtocol, = struct.unpack('<i', _buffer)
    print(f'Header=>NetworkProtocol: {networkProtocol}')
    # 04 Server Name
    _buffer = iFile.read(260) # 260 defined by valve
    serverName = _buffer.decode()
    print(f'Header=>ServerName: {serverName}')
    # 05 Client Name
    _buffer = iFile.read(260) # 260 defined by valve
    clientName = _buffer.decode()
    print(f'Header=>ClientName: {clientName}')
    # 06 Map Name
    _buffer = iFile.read(260)
    mapName = _buffer.decode()
    print(f'Header=>MapName: {mapName}')
    # 07 Game Directory
    _buffer = iFile.read(260)
    gameDir = _buffer.decode()
    print(f'Header=>GameDirectory: {gameDir}')
    # 08 Playback Time [seconds]
    _buffer = iFile.read(4)
    playTime, = struct.unpack('f', _buffer)
    print(f'Header=>PlaybackTime: {playTime}')
    # 09 Ticks
    _buffer = iFile.read(4)
    ticks, = struct.unpack('i', _buffer)
    print(f'Header=>Ticks: {ticks}')
    # 10 Frames
    _buffer = iFile.read(4)
    frames, = struct.unpack('i', _buffer)
    print(f'Header=>Frames: {frames}')
    # 11 Sign on length [for data signon] ??
    _buffer = iFile.read(4)
    signon, = struct.unpack('i', _buffer)
    print(f'Header=>Sign Length: {signon}')
    
    tickrate = ticks / playTime
    print(f'tickrate: {tickrate}')




def parseDemoFile(iFile: io.BufferedReader):
    # Demo Header
    parseHeader(iFile)
    _buffer = iFile.read(4)
    command, = struct.unpack('i', _buffer)
    print(command)
    pass

def main():
    with open(demofile, 'rb') as iFile:
        parseDemoFile(iFile)

if __name__ == '__main__':
    main()