package fr.membrives.etienne.ringlistener;

import android.bluetooth.BluetoothDevice;
import android.bluetooth.BluetoothGattCharacteristic;
import android.util.Log;

import java.util.Calendar;

/**
 * Created by etienne on 02/04/16.
 */
public class Controller implements BluetoothConnect.Callback {
    private static String TAG = "Controller";

    @Override
    public void onConnected(BluetoothConnect uart) {
        Log.d(TAG, "onConnected");
    }

    @Override
    public void onConnectFailed(BluetoothConnect uart) {
        Log.d(TAG, "onConnectFailed");
    }

    @Override
    public void onDisconnected(BluetoothConnect uart) {
        Log.d(TAG, "onDisconnected");
    }

    @Override
    public void onReceive(BluetoothConnect uart, BluetoothGattCharacteristic rx) {
        String command = rx.getStringValue(0);
        Log.d(TAG, "Command: " + command);
        if (command.equals("ring")) {

        } else if (command.equals("time")) {
            Calendar c = Calendar.getInstance();
            int hours = c.get(Calendar.HOUR_OF_DAY);
            int minutes = c.get(Calendar.MINUTE);
            uart.send(String.format("%d:%d", hours, minutes));
        }
    }

    @Override
    public void onDeviceFound(BluetoothDevice device) {
        Log.d(TAG, "onDeviceFound");
    }

    @Override
    public void onDeviceInfoAvailable() {
        Log.d(TAG, "onDeviceInfoAvailable");
    }
}
