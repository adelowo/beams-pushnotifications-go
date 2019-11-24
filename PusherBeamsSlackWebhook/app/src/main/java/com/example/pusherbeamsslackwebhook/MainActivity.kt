package com.example.pusherbeamsslackwebhook

import android.content.SharedPreferences
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import com.pusher.pushnotifications.*
import com.pusher.pushnotifications.auth.AuthData
import com.pusher.pushnotifications.auth.AuthDataGetter
import com.pusher.pushnotifications.auth.BeamsTokenProvider
import java.util.*


class MainActivity : AppCompatActivity() {

    private val PREF_NAME = "uuid-generated"


    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)


        PushNotifications.start(applicationContext, "PUSHER_INSTANCE_ID")
        PushNotifications.addDeviceInterest("webhook-slack")

        val sharedPref: SharedPreferences = getSharedPreferences(PREF_NAME, 0)


        if (!sharedPref.getBoolean(PREF_NAME, false)) {

            var uuid = UUID.randomUUID().toString()

            // get the token from the server
            val serverUrl = "https://NGROK_URL/auth?user_id=${uuid}"
            val tokenProvider = BeamsTokenProvider(serverUrl,
                object : AuthDataGetter {
                    override fun getAuthData(): AuthData {
                        return AuthData(
                            headers = hashMapOf()
                        )
                    }
                })


            PushNotifications.setUserId(
                uuid,
                tokenProvider,
                object : BeamsCallback<Void, PusherCallbackError> {
                    override fun onFailure(error: PusherCallbackError) {
                        Log.e(
                            "BeamsAuth",
                            "Could not login to Beams: ${error.message}"
                        )
                    }

                    override fun onSuccess(vararg values: Void) {
                        Log.i("BeamsAuth", "Beams login success")
                    }
                }
            )
            val editor = sharedPref.edit()
            editor.putBoolean(PREF_NAME, true)
            editor.apply()
        }
    }
}
