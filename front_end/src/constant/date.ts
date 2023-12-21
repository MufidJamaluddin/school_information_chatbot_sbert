const dateIDFormatter = new Intl.DateTimeFormat(
  'id-ID',  
  { 
    dateStyle: 'full', 
    timeStyle: 'long', 
    timeZone: 'Asia/Jakarta' 
  }
);

export {
  dateIDFormatter
}