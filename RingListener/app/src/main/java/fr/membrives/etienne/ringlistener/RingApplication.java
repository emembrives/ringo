package fr.membrives.etienne.ringlistener;

import android.app.Application;
import android.content.Intent;
import android.os.RemoteException;
import android.util.Log;

import org.altbeacon.beacon.Beacon;
import org.altbeacon.beacon.BeaconConsumer;
import org.altbeacon.beacon.BeaconManager;
import org.altbeacon.beacon.BeaconParser;
import org.altbeacon.beacon.Identifier;
import org.altbeacon.beacon.RangeNotifier;
import org.altbeacon.beacon.Region;
import org.altbeacon.beacon.powersave.BackgroundPowerSaver;
import org.altbeacon.beacon.startup.BootstrapNotifier;
import org.altbeacon.beacon.startup.RegionBootstrap;
import org.altbeacon.beacon.utils.UrlBeaconUrlCompressor;

import java.util.Collection;
import java.util.logging.LogManager;

/**
 * Created by etienne on 20/02/16.
 */
public class RingApplication extends Application implements BootstrapNotifier, BeaconConsumer,
        RangeNotifier {
    private static final String TAG = ".RingApplication";
    private RegionBootstrap regionBootstrap;
    private BeaconManager beaconManager;
    private BackgroundPowerSaver backgroundPowerSaver;
    private Region mRegion = new Region("fr.membrives.ringlistener.ringer",
            Identifier.parse("0x0265746e2e6f76682f72"), null,
            null);
    private BluetoothConnect mBluetoothConnect;
    private final Controller mController = new Controller();

    @Override
    public void onCreate() {
        super.onCreate();
        Log.d(TAG, "App started up");
        beaconManager = BeaconManager.getInstanceForApplication(this);
        // To detect proprietary beacons, you must add a line like below corresponding to your beacon
        // type.  Do a web search for "setBeaconLayout" to get the proper expression.
        // beaconManager.getBeaconParsers().add(new BeaconParser().
        //        setBeaconLayout("m:2-3=beac,i:4-19,i:20-21,i:22-23,p:24-24,d:25-25"));

        // Eddystone beacons
        // Detect the main identifier (UID) frame:
        beaconManager.getBeaconParsers().add(new BeaconParser().
                setBeaconLayout("s:0-1=feaa,m:2-2=00,p:3-3:-41,i:4-13,i:14-19"));
        // Detect the telemetry (TLM) frame:
        beaconManager.getBeaconParsers().add(new BeaconParser().
                setBeaconLayout("x,s:0-1=feaa,m:2-2=20,d:3-3,d:4-5,d:6-7,d:8-11,d:12-15"));
        // Detect the URL frame:
        beaconManager.getBeaconParsers().add(new BeaconParser().
               setBeaconLayout("s:0-1=feaa,m:2-2=10,p:3-3:-41,i:4-20v"));
        beaconManager.bind(this);

        // wake up the app when any beacon is seen (you can specify specific id filers in the parameters below)
        regionBootstrap = new RegionBootstrap(this, mRegion);
        backgroundPowerSaver = new BackgroundPowerSaver(this);

        mBluetoothConnect = new BluetoothConnect(this);
        mBluetoothConnect.registerCallback(mController);
    }

    @Override
    public void didDetermineStateForRegion(int arg0, Region arg1) {
        Log.d(TAG, "Got a didDetermineStateForRegion call");
        // Don't care
    }

    @Override
    public void didEnterRegion(Region region) {
        Log.d(TAG, "Got a didEnterRegion call");
        try {
            beaconManager.startRangingBeaconsInRegion(mRegion);
        } catch (RemoteException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void didExitRegion(Region arg0) {
        Log.d(TAG, "Got a didExitRegion call");
        try {
            beaconManager.stopRangingBeaconsInRegion(mRegion);
        } catch (RemoteException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void onBeaconServiceConnect() {
        Log.d(TAG, "Got a onBeaconServiceConnect call");
        beaconManager.setRangeNotifier(this);
    }

    @Override
    public void didRangeBeaconsInRegion(Collection<Beacon> beacons, Region region) {
        Log.d(TAG, "Got a didRangeBeaconsInRegion call");
        for (Beacon beacon: beacons) {
            if (beacon.getServiceUuid() == 0xfeaa && beacon.getBeaconTypeCode() == 0x10) {
                // This is a Eddystone-URL frame
                String url = UrlBeaconUrlCompressor.uncompress(beacon.getId1().toByteArray());
                Log.d(TAG, "I see a beacon transmitting a url: " + url +
                           " approximately " + beacon.getDistance() + " meters away.");
                Log.d(TAG, "With id " + beacon.getId1().toString());
                try {
                    beaconManager.stopRangingBeaconsInRegion(mRegion);
                } catch (RemoteException e) {
                    e.printStackTrace();
                }
                mBluetoothConnect.connect(beacon);
            }
        }
    }
}
