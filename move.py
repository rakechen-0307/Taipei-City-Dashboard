import os
from time import sleep

while True:
	if os.path.isfile("./Taipei-City-Dashboard-BE/incident.geojson"):
		print("file exist")
		os.system("rm -f ./Taipei-City-Dashboard-FE/public/mapData/incident.geojson")
		os.system("mv ./Taipei-City-Dashboard-BE/incident.geojson ./Taipei-City-Dashboard-FE/public/mapData/incident.geojson")
	sleep(10)