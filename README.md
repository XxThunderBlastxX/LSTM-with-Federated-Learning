# 🤖 LSTM-with-Federated-Learning

## 🛠️ Pre-requisites
1. [GoLang](https://go.dev/dl/)
2. [Python v3.10.13](https://www.python.org/downloads/release/python-31013/)


## ⚙️ How to run the code
1. Install 64bit Python v3.10.13

2. Install the required packages using the following command:
```bash
pip install -r requirements.txt
```

3. Run server using the following command:
```bash
go run main.go
```

5. Run Client1.py, Client2.py, Client3.py -> it will output the loss, RMSE, NRMSE for each client.

```bash
python Client<client_number>.py
```

6. Currently the stop criteria has not be been set up. It will continue to run. Final accuracy will be approaching 1 eventually and weights converge (weights do not change at local client and the weights of two client are close to each other).
