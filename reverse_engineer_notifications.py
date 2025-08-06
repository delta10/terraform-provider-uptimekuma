#!/usr/bin/env python3
"""
Reverse engineer Uptime Kuma's notification API using Socket.IO
This script will help us understand the proper way to create and manage notifications
"""

import socketio
import json
import time
import sys

# Create a Socket.IO client
sio = socketio.Client()

# Global variables to store received data
received_events = []
notifications_data = None
monitor_list_data = None

@sio.event
def connect():
    print("Connected to Uptime Kuma Socket.IO server")

@sio.event
def disconnect():
    print("Disconnected from server")

@sio.event
def connect_error(data):
    print(f"Connection failed: {data}")

# Catch all events to see what Uptime Kuma sends
@sio.on('*')
def catch_all(event, *args):
    global notifications_data, monitor_list_data
    print(f"Received event: {event}")
    print(f"Args: {json.dumps(args, indent=2, default=str)}")
    print("-" * 50)
    
    received_events.append({'event': event, 'args': args})
    
    # Store specific data we're interested in
    if event == 'notificationList':
        notifications_data = args[0] if args else None
        print(f"NOTIFICATION LIST DATA: {json.dumps(notifications_data, indent=2, default=str)}")
    elif event == 'monitorList':
        monitor_list_data = args[0] if args else None

def login(username, password):
    """Login to Uptime Kuma"""
    print(f"Attempting to login as {username}...")
    
    def login_callback(response):
        print(f"Login response: {json.dumps(response, indent=2, default=str)}")
        if response.get('ok'):
            print("✅ Login successful!")
            return True
        else:
            print("❌ Login failed!")
            return False
    
    sio.emit('login', {
        'username': username,
        'password': password,
        'token': ''
    }, callback=login_callback)
    
    time.sleep(2)  # Wait for response

def get_notification_list():
    """Try to get the notification list"""
    print("Requesting notification list...")
    
    def callback(response):
        print(f"getNotificationList response: {json.dumps(response, indent=2, default=str)}")
    
    sio.emit('getNotificationList', {}, callback=callback)
    time.sleep(2)

def create_test_notification():
    """Try to create a test notification"""
    print("Creating test notification...")
    
    test_notification = {
        'name': 'Socket.IO Test Notification',
        'type': 'webhook',
        'isDefault': False,
        'applyExisting': False,
        'webhookURL': 'https://httpbin.org/post',
        'webhookContentType': 'application/json'
    }
    
    def callback(response):
        print(f"addNotification response: {json.dumps(response, indent=2, default=str)}")
    
    sio.emit('addNotification', test_notification, callback=callback)
    time.sleep(3)

def main():
    try:
        # Connect to Uptime Kuma
        print("Connecting to Uptime Kuma at http://localhost:3001...")
        sio.connect('http://localhost:3001', socketio_path='/socket.io/')
        
        # Wait for initial connection
        time.sleep(1)
        
        # Login
        login('admin', 'cF96H*L9LA3*HiWhx')
        
        # Wait for any initial events
        time.sleep(2)
        
        # Try to get notification list
        get_notification_list()
        
        # Try different notification list calls
        print("Trying 'getNotifications'...")
        sio.emit('getNotifications', {}, callback=lambda x: print(f"getNotifications: {x}"))
        time.sleep(2)
        
        print("Trying 'getSettings'...")
        sio.emit('getSettings', {}, callback=lambda x: print(f"getSettings: {x}"))
        time.sleep(2)
        
        # Create a test notification
        create_test_notification()
        
        # Try to get notification list again after creating
        get_notification_list()
        
        # Keep connection alive to receive any delayed events
        print("Waiting for any additional events...")
        time.sleep(5)
        
        print("\n" + "="*60)
        print("SUMMARY OF RECEIVED EVENTS:")
        print("="*60)
        for event in received_events:
            print(f"Event: {event['event']}")
        
        if notifications_data:
            print(f"\nNotification List Data Structure:")
            print(json.dumps(notifications_data, indent=2, default=str))
        
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()
    finally:
        sio.disconnect()

if __name__ == "__main__":
    main()
