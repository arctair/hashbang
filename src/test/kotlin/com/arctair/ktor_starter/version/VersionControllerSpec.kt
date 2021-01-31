package com.arctair.ktor_starter.version

import com.arctair.ktor_starter.globalModules
import io.ktor.http.*
import io.ktor.server.testing.*
import io.mockk.every
import io.mockk.impl.annotations.MockK
import io.mockk.junit5.MockKExtension
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.json.Json
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.extension.ExtendWith
import kotlin.test.assertEquals

@ExtendWith(MockKExtension::class)
internal class VersionControllerSpec {

  @MockK
  lateinit var versionService: VersionService

  private fun <R> withController(test: TestApplicationEngine.() -> R) =
    withTestApplication(
      {
        versionController(versionService)
        globalModules()
      },
      test,
    )

  @Test
  fun `get version`() = withController {
    val givenVersion = Version("the sha", "the version")
    every { versionService.getVersion() } returns givenVersion
    val call = handleRequest(method = HttpMethod.Get, uri = "/version")

    assertEquals(HttpStatusCode.OK, call.response.status())

    val version = Json.decodeFromString<Version>(call.response.content!!)
    assertEquals(givenVersion, version)
  }
}