package example

type RezMessage struct {
    
        Dummy bool
    
    Dirty bool
    cache *example.Rez
}

func NewRezMessageFromFbs(fbs *example.Rez) (m *RezMessage) {
    return &RezMessage{
         
            Dummy: fbs.Dummy(),
         
         Dirty: false,
         cache: fbs,   
    }
}

func (m *RezMessage) ApplyUpdated() bool {
    
        if m.Dummy != m.cache.Dummy() {
             return m.cache.MutateDummy(m.Dummy)
            
            // if hasnt mutate method, return false 
            //  return false
        }
    

    return true
}

func (m *RezMessage) ToFbs(check bool) (fbs *example.Rez) {
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
