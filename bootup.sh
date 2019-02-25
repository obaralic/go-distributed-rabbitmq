#!/bin/bash

# cd coordinator/exec/
# go clean
# go build exec.go
# gnome-terminal -x exec
# cd ../..
#
# cd web/
# go clean
# go build exec.go
# gnome-terminal -x exec
# cd ..

cd sensors/
go clean
go build sensors.go
gnome-terminal -x sensors -name=boiler_pressure_out    -min=15    -max=15.5  -dev=0.05   -freq=1
# gnome-terminal -x sensors -name=turbine_pressure_out   -min=0.9   -max=1.3   -dev=0.05   -freq=1
# gnome-terminal -x sensors -name=condensor_pressure_out -min=0.001 -max=0.002 -dev=0.0001 -freq=1
# gnome-terminal -x sensors -name=boiler_temp_out        -min=590   -max=615   -dev=1      -freq=1
# gnome-terminal -x sensors -name=turbine_temp_out       -min=100   -max=105   -dev=1      -freq=1
# gnome-terminal -x sensors -name=condensor_temp_out     -min=80    -max=98    -dev=1      -freq=1
cd ..
