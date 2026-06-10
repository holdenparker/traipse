# Traipse

The goal of this project is to gain experience in golang while building out OGC compliant services.
This is a hobby project at best and will not receive regular maintenance.
_**This is NOT to be utilized by any production system.**_

## geo-types

Golang structs that map from GeoJSON features and geometries.
Eventually I would like to explore GML and KML and other related types.

## algorithms

Holds spatial algorithms to be utilized by the rest of the system.

## ogc-api

Where I am going to be building out OGC compliant APIs (hopefully).

### features

I am taking on the OGC API Features Part 1 Core service first.
Once I've had enough of the OGC API Features service, I may move on to Tiles or something else that I find interesting at the time.

#### schemas

In an attempt to streamline my conformance with the OGC API spec, I have pulled in OAPI yaml and utilized that to generate structs for use in the Huma API.
I intentionally ignored the GeoJSON-related types, even though they are very similar to the ones I've created, simply because I have other plans for the `geo-types` structs.