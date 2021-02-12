package com.arctair.hashbang.bspec

import io.ktor.client.*
import io.ktor.client.features.*
import io.ktor.client.features.json.*
import io.ktor.client.request.*
import io.ktor.http.*
import org.junit.jupiter.api.extension.BeforeEachCallback
import org.junit.jupiter.api.extension.ExtensionContext

@Target(AnnotationTarget.FIELD)
internal annotation class TestHttpClient(val clientName: String = "")

internal object TestHttpClientProvider : BeforeEachCallback, ExtensionContext.Store.CloseableResource {

  private val httpClientByName = mapOf(
    "" to HttpClient {
      install(JsonFeature)
      defaultRequest {
        setBaseUrl(this, "targetBaseUri", "http://localhost:8080")
      }
    },
  )

  private fun setBaseUrl(builder: HttpRequestBuilder, name: String, default: String) {
    val baseUri = System.getProperty(name) ?: default
    val (protocolString, hostname, basePath) = Regex("(?<protocol>\\w+)://(?<hostname>[\\w.-]+(:\\d+)?)(?<basePath>.*)")
      .matchEntire(baseUri)
      ?.groupValues
      ?.slice(listOf(1, 2, 4)) ?: throw Error("URL environment property $name='$baseUri' is malformed")
    builder.url.protocol = URLProtocol.createOrDefault(protocolString)
    builder.url.host = hostname
    builder.url.encodedPath =
      if (builder.url.encodedPath.isEmpty()) basePath
      else "${basePath}/${builder.url.encodedPath}"
  }

  override fun beforeEach(context: ExtensionContext?) {
    val instance = context!!.requiredTestInstance
    instance.javaClass.declaredFields
      .filter { it.isAnnotationPresent(TestHttpClient::class.java) }
      .forEach { it.set(instance, httpClientByName[it.getAnnotation(TestHttpClient::class.java).clientName]) }
  }

  override fun close() {
    httpClientByName.values.forEach(HttpClient::close)
  }
}
