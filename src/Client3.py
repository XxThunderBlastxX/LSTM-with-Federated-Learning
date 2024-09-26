import websocket
import json
import numpy as np
from TFClass3 import TFML
import sys

def on_message(ws, message):
    global clientModel, new_client_weight
    print("Received message from server:", message[:100] + "..." if len(message) > 100 else message)
    try:
        received_weights = json.loads(message)
        
        # Convert the received weights back to numpy arrays
        received_weights = [np.array(w) for w in received_weights]
        
        # Print received weight shapes
        print("Received weight shapes:", [w.shape for w in received_weights])
        
        # Get current model weights
        current_weights = clientModel.model.get_weights()
        print("Current weight shapes:", [w.shape for w in current_weights])
        
        # Update weights where possible
        for i, (received, current) in enumerate(zip(received_weights, current_weights)):
            if received.shape == current.shape:
                current_weights[i] = received
            elif received.size <= current.size:
                flat_current = current.flatten()
                flat_current[:received.size] = received.flatten()
                current_weights[i] = flat_current.reshape(current.shape)
            else:
                print(f"Warning: Received weight {i} is larger than current weight. Skipping update.")
        
        # Set the updated weights
        clientModel.model.set_weights(current_weights)
        print("Weights updated. New shapes:", [w.shape for w in clientModel.model.get_weights()])
        
        print("Running model...")
        clientModel.run()
        print("Evaluating model...")
        clientModel.eval()
        new_client_weight = clientModel.model.get_weights()
        
        # Send the new weights back to the server
        send_weights(ws)
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
    except Exception as e:
        print(f"Unexpected error in on_message: {e}")
        import traceback
        traceback.print_exc()


def on_error(ws, error):
    print(f"WebSocket error: {error}")

def on_close(ws, close_status_code, close_msg):
    print(f"Connection closed. Status code: {close_status_code}. Message: {close_msg}")

def on_open(ws):
    print("Connection established")
    send_weights(ws)

def send_weights(ws):
    global clientModel, new_client_weight
    print(f"Preparing to send weights to server (Client ID: {clientModel.name})")
    try:
        weights_to_send = [w.tolist() for w in new_client_weight]  # Convert numpy arrays to lists
        print(f"Weight structure: {[np.array(w).shape for w in weights_to_send]}")  # Print weight shapes
        data = {
            "weights": weights_to_send,
            "clientId": clientModel.name
        }
        json_data = json.dumps(data)
        print(f"Sending data of size: {sys.getsizeof(json_data)} bytes")
        ws.send(json_data)
        print("Weights sent successfully")
    except Exception as e:
        print(f"Error sending weights: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    print("Initializing TFML model...")
    clientModel = TFML('client3')
    new_client_weight = clientModel.model.get_weights()

    print('Waiting for connection...')
    websocket.enableTrace(True)  # Enable debug trace
    ws = websocket.WebSocketApp("ws://127.0.0.1:1233/ws",
                                on_message=on_message,
                                on_error=on_error,
                                on_close=on_close,
                                on_open=on_open)
    
    print("Starting WebSocket connection...")
    ws.run_forever()