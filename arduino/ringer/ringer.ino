/*********************************************************************
 This is an example for our nRF51822 based Bluefruit LE modules

 Pick one up today in the adafruit shop!

 Adafruit invests time and resources providing this open source code,
 please support Adafruit and open-source hardware by purchasing
 products from Adafruit!

 MIT license, check LICENSE for more information
 All text above, and the splash screen below must be included in
 any redistribution
*********************************************************************/

#include <Arduino.h>
#include <SPI.h>
#if not defined (_VARIANT_ARDUINO_DUE_X_) && not defined (_VARIANT_ARDUINO_ZERO_)
  #include <SoftwareSerial.h>
#endif

#include "Adafruit_BLE.h"
#include "Adafruit_BluefruitLE_SPI.h"
#include "Adafruit_BluefruitLE_UART.h"

#include "BluefruitConfig.h"

/*=========================================================================
    APPLICATION SETTINGS

    FACTORYRESET_ENABLE       Perform a factory reset when running this sketch

                              Enabling this will put your Bluefruit LE module
                              in a 'known good' state and clear any config
                              data set in previous sketches or projects, so
                              running this at least once is a good idea.

                              When deploying your project, however, you will
                              want to disable factory reset by setting this
                              value to 0.  If you are making changes to your
                              Bluefruit LE device via AT commands, and those
                              changes aren't persisting across resets, this
                              is the reason why.  Factory reset will erase
                              the non-volatile memory where config data is
                              stored, setting it back to factory default
                              values.

                              Some sketches that require you to bond to a
                              central device (HID mouse, keyboard, etc.)
                              won't work at all with this feature enabled
                              since the factory reset will clear all of the
                              bonding data stored on the chip, meaning the
                              central device won't be able to reconnect.
    MINIMUM_FIRMWARE_VERSION  Minimum firmware version to have some new features
    URL                       The URL that is advertised. It must not longer
                              than 17 bytes (excluding http:// and www.).
                              Note: ".com/" ".net/" count as 1
    -----------------------------------------------------------------------*/
    #define FACTORYRESET_ENABLE         1
    #define MINIMUM_FIRMWARE_VERSION    "0.6.7"
    #define URL                         "http://etn.ovh/r"
/*=========================================================================*/


// Create the bluefruit object, with hardware SPI, using SCK/MOSI/MISO
// hardware SPI pins and then user selected CS/IRQ/RST
Adafruit_BluefruitLE_SPI ble(BLUEFRUIT_SPI_CS, BLUEFRUIT_SPI_IRQ, BLUEFRUIT_SPI_RST);

// A small helper
void error(const __FlashStringHelper*err) {
  while (!Serial);  // required for Flora & Micro
  Serial.begin(115200);
  Serial.println(err);
  while (1);
}

/**************************************************************************/
/*!
    @brief  Sets up the HW an the BLE module (this function is called
            automatically on startup)
*/
/**************************************************************************/
void setup(void)
{
  delay(500);

  if ( !ble.begin(VERBOSE_MODE) )
  {
    error(F("Couldn't find Bluefruit, make sure it's in CoMmanD mode & check wiring?"));
  }

  if ( FACTORYRESET_ENABLE )
  {
    /* Perform a factory reset to make sure everything is in a known state */
    if ( ! ble.factoryReset() ){
      error(F("Couldn't factory reset"));
    }
  }

  /* Disable command echo from Bluefruit */
  ble.echo(false);

  // EddyStone commands are added from firmware 0.6.6
  if (!ble.isVersionAtLeast(MINIMUM_FIRMWARE_VERSION) )
  {
    error(F("EddyStone is only available from 0.6.6. Please perform firmware upgrade"));
  }

  if (!ble.sendCommandCheckOK(F( "AT+EDDYSTONEURL=" URL ))) {
    error(F("Couldnt set, is URL too long !?"));
  }

  if (!ble.sendCommandCheckOK(F("AT+EDDYSTONEENABLE=on")) ) {
    error(F("Couldnt enable Eddystone"));
  }

  ble.sendCommandCheckOK(F("AT+EDDYSTONEENABLE"));
}

/**************************************************************************/
/*!
    @brief  Constantly poll for new command or response data
*/
/**************************************************************************/
void loop(void)
{
  delay(1000);
  if (ble.isConnected()) {
    ble.println("AT+BLEUARTRX");
    ble.readline();
    
    if (strcmp(ble.buffer, "OK") != 0) {
      // Some data was found, its in the buffer
      Serial.print(F("[Recv] "));
      Serial.println(ble.buffer);
      if (!ble.waitForOK()) {
        error(F("Failed to receive"));
      }
    }
  
    ble.print("AT+BLEUARTTX=");
    ble.println("ring");
    if (!ble.waitForOK()) {
      error(F("Failed to send"));
    }
  }
}

