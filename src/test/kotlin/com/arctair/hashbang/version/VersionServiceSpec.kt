package com.arctair.hashbang.version

import org.hamcrest.Matchers
import org.hamcrest.junit.MatcherAssert
import org.junit.jupiter.api.Test

internal class VersionServiceSpec {
  
  private val versionService = VersionService()

  @Test
  fun `get version`() {
      val version = versionService.getVersion()

      MatcherAssert.assertThat(version.sha, Matchers.matchesPattern("[0-9a-f]{40}"))
      MatcherAssert.assertThat(version.version, Matchers.matchesPattern("v\\d+\\.\\d+\\.\\d+"))
  }
}