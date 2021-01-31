package com.arctair.ktor_starter.bspec

import com.arctair.ktor_starter.applicationModule
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import org.junit.jupiter.api.extension.BeforeAllCallback
import org.junit.jupiter.api.extension.ExtensionContext

internal object StartApplication : BeforeAllCallback, ExtensionContext.Store.CloseableResource {
  var server: NettyApplicationEngine? = null

  override fun beforeAll(context: ExtensionContext?) {
    if (System.getProperty("ktor_starter_baseUrl") != null || server != null) return
    server = embeddedServer(
      Netty,
      applicationEngineEnvironment {
        module {
          applicationModule()
        }
        connector {
          port = 8080
        }
      },
    ).start()
  }

  override fun close() {
    server?.stop(1000, 10000)
  }
}