package com.leg3nd.plugins

import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.plugins.statuspages.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.serialization.Serializable

fun Application.configureRouting() {
    install(StatusPages) {
        exception<Throwable> { call, cause ->
            call.respondText(text = "500: $cause", status = HttpStatusCode.InternalServerError)
        }
    }
    routing {
        get("/") {
            val lb = listOf(Foo.Bar("hi"), Foo.Bar("my name is"), Foo.Bar("d0lim"))
            val f = Foo(lb)
            call.respond(f)
        }
    }
}

@Serializable
data class Foo(
    val some: List<Bar>
) {
    @Serializable
    data class Bar(val name: String)
}