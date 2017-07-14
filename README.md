# arcrestgo
Go version of simple rest server

docker run -it -v /d/data:/data -v /d/git:/git geodata/gdal:local bash

docker run -it --name pg -p 5432:5432 -v /d/data:/data -v /d/git:/git kartoza/postgis:9.4-2.1

ogrinfo -ro PG:'dbname=gis host=192.168.99.100 user=postgres' -sql "SELECT count(*) from spatial_ref_sys"

Load data from sqlite to postgresql

ogr2ogr -f "PostgreSQL" PG:"host=192.168.99.100 dbname=gis user=postgres" /data/accommodationagreementrentals.geodatabase

ogr2ogr -append -lco GEOMETRY_NAME=the_geom -lco SCHEMA=public -f "PostgreSQL" PG:"host= port=5432 user=postgres dbname=gis" -a_srs "EPSG:3857" accommodationagreementrentals.geodatabase

ogr2ogr -f "PostgreSQL" PG:"dbname=gis user=postgres" "source_data.json" 


docker run  -v /d/data:/data  -v /d/git:/git geodata/gdal:local ogr2ogr -f "PostgreSQL" PG:"host=192.168.99.100 dbname=gis user=postgres" /data/accommodationagreementrentals.geodatabase -skipfailures
docker run  -v /d/data:/data  -v /d/git:/git geodata/gdal:local ogr2ogr -f "PostgreSQL" PG:"host=192.168.99.100 dbname=gis user=postgres" /data/content.items.data.json -skipfailures


SELECT load_extension( 'd:/bin/mod_spatialite.dll', 'sqlite3_modspatialite_init')
SELECT load_extension( 'd:/bin/stgeometry_sqlite.dll', 'SDE_SQL_funcs_init');

Docker
git clone git@github.com:traderboy/arcrestgo.git
cd arcrestgo
cd docker
docker-compose build
docker-compose up -d