<?xml version="1.0" encoding="UTF-8"?>
<tileset version="1.10" tiledversion="1.11.2" name="atlas" tilewidth="16" tileheight="16" tilecount="336" columns="21">
 <image source="atlas.png" width="336" height="256"/>
 <tile id="7">
  <objectgroup draworder="index" id="2">
   <object id="1" x="16" y="7">
    <polygon points="0,0 -1,0 -7,3 -8,4 -10,8 -10,9 0,9"/>
   </object>
   <object id="2" x="8" y="0" width="8" height="8"/>
   <object id="3" x="6" y="2" width="5" height="9" rotation="34.6505">
    <ellipse/>
   </object>
  </objectgroup>
 </tile>
 <tile id="8">
  <objectgroup draworder="index" id="2">
   <object id="1" x="6" y="6">
    <polygon points="0,0 -2,1 -6,1 -6,10 10,10 10,1 5,1 4,0"/>
   </object>
  </objectgroup>
 </tile>
 <tile id="110">
  <objectgroup draworder="index" id="2">
   <object id="1" x="5" y="4">
    <point/>
   </object>
  </objectgroup>
 </tile>
 <tile id="111">
  <properties>
   <property name="test" value="hello, world!"/>
  </properties>
 </tile>
 <tile id="115">
  <objectgroup draworder="index" id="2">
   <object id="4" name="collision" x="9" y="4">
    <properties>
     <property name="solid" type="bool" value="true"/>
    </properties>
    <polyline points="1,7 -1,12 7,12 7,5 1,7"/>
   </object>
   <object id="5" name="origin" x="7" y="6">
    <point/>
   </object>
  </objectgroup>
 </tile>
 <tile id="138">
  <objectgroup draworder="index" id="2">
   <object id="1" x="0" y="0" width="9" height="16"/>
  </objectgroup>
 </tile>
 <tile id="139">
  <animation>
   <frame tileid="110" duration="100"/>
   <frame tileid="111" duration="100"/>
   <frame tileid="112" duration="200"/>
   <frame tileid="133" duration="100"/>
   <frame tileid="154" duration="100"/>
   <frame tileid="153" duration="100"/>
   <frame tileid="152" duration="100"/>
   <frame tileid="131" duration="100"/>
  </animation>
 </tile>
</tileset>
