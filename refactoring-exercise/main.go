package router

// ComputeVisHandler computes the satellite visibility
func(s SatelliteRouter) ComputeVisHandler(start, stop time.Time, ephemerisType, fileName string) map[string][] string {
    // queries the DB for all satellites
    satellites, err: = db.ListAllSatellites()
    // stores the sat name to DB ID
    s: = make(map[string] int64, 0)
    // stores the gateway names to DB ID
    g: = make(map[string] int64, 0)
    for _, ss: = range satellites {
            s[ss.Name] = ss.ID
        }

	/// MAX COMMENT: the operation is exacly the same for Satellites
	/// and Gateways.
	/// Consider to: 
	/// is the re-indexing by name useful?
	/// can't make ListAllSatellites to return a hashmap?
	/// maybe re-indexing by name through a function?
	/// maybe is better to introduce a dedicated type struct for our satellite?
	/// END COMMENT

    // queries the DB for all gateways
    gateways, err: = db.ListAllGateways()
    for _, gg: = range gateways {
            g[gg.Name] = gg.ID
        }
    // queries the DB for all satellite parameters
    sp, err: = db.ListAllsatelliteParams()

	// MAX COMMENT:
	/// really can't this be done while querying for Satellites?
	// END COMMENT

	// the following could be a computationally long function.
	/// MAX COMMENT:
	/// maybe is better to queue the operation and execute it in 
	/// asyncronus way, if the effect of the function is not required to
	/// be immediate and blocking?
	/// END COMMENT
    err = ArchiveFileToAllAgents(fileName)

    v: = make(map[string][] string)
    if fileType == "TLE" {
        for _, ss: = sp {  /// MAX COMMENT: No range here?
            // computes satellite visibility base on satellite properties
            v: = computeVisTle(ss, s, g, start, stop)
        }
    } ///  MAX COMMENT: An else would be better
    if fileType == "OEM" {
        for _, ss: = sp { /// MAX COMMENT: No range here?
            // computes satellite visibility base on satellite properties
            v: = computeVisOem(ss, s, g, start, stop)
        }
    } ///  MAX COMMENT: Need an else with an error raised

	///  MAX COMMENT: Only the function change
	/// this can be rewritten in a cooler way
	if fileType == "TLE" {
		compute_func := computeVisTle
	} else {
		compute_func := computeVisOem
	}
	v: = compute_func(ss, s, g, start, stop)
	/// END COMMENT

    return v
}