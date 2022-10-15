# MechAssist

The goal for this project is a simple web app that solves for the cycle efficiency of steam power cycles. The advantage it offers to mechanical students is that you see the state of steam at the different processes and then you know where your calculation went wrong. 


## Getting Started
    1. git clone this repo
    2. must have go 1.19 installed on your computer
    3. go run main.go!

### Sample responses
    The newRankine_sat method takes two pressure units, pressure at the condenser and pressure at the boiler respectively.
    The sat keyword stands for saturated i.e saturated vapor entering the turbine. Currently working on functions for when superheated steam enters the turbine.
```
    newRankine_sat(0.01, 2)

    +-----------+---------------+-----------------------------+-----------------------------+
    |  PROCESS  |  INLET STATE  |         EXIT STATE          |   WORK DONE/HEAT TRANSFER   |
    +-----------+---------------+-----------------------------+-----------------------------+
    | Pump      | P1 = 0.01MPa  | h2 = 193.83MPa              | Work Input: 2.0153KJ/KG     |
    | Boiler    | P2 = 2MPa     | h3 = 2798.3MPa              | Heat Added : 2604.5KJ/KG    |
    | Turbine   | State 3 known | x4 = 0.75868MPa h4 = 2006.6 | Work output : 791.65KJ/KG   |
    | Condenser | State 4 known | Output known                | Heat rejected : 1814.8KJ/KG |
    +-----------+---------------+-----------------------------+-----------------------------+
    |                                 EFFICIENCY OF CYCLE     |           30.318%           |
    +-----------+---------------+-----------------------------+-----------------------------+
```
Contributions:
    If you are a mechanical engineering or Chemistry student interested in extending the functionality, fork and contribute. Let's make something great!
