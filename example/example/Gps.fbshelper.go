package example

type GpsMessage struct {
    
        Latitude float64
    
        Longitude float64
    
    Dirty bool
    cache *example.Gps
}

func NewGpsMessageFromFbs(fbs *example.Gps) (m *GpsMessage) {
    return &GpsMessage{
         
            Latitude: fbs.Latitude(),
         
            Longitude: fbs.Longitude(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *GpsMessage) ApplyUpdated() bool {
    
        if m.Latitude != m.cache.Latitude() {
             return m.cache.MutateLatitude(m.Latitude)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    
        if m.Longitude != m.cache.Longitude() {
             return m.cache.MutateLongitude(m.Longitude)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *GpsMessage) ToFbs(check bool) (fbs *example.Gps) {
    if !check && m.cache != nil {
        return m.cache
    }
    if m.ApplyUpdated() {
        return m.cache
    }
    m.Dirty = true
    // 作成処理
    /*
          fbs = ...   
          m.cache = fbs
    */

    return m.cache
}
