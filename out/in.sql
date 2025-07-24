select s.nda,p.sexe,to_number(substr(p.DATENAIS  ,1, 4)) AS anneenais, r.nir,l.nith, reponse, srv.nom AS service, niqsup,nirsup,l.parente,
CASE
	WHEN VALIDATE_CONVERSION(reponse AS DATE, 'DD/MM/YYYY') = 1
		THEN TO_DATE(reponse,'DD/MM/YYYY')
	ELSE NULL
END AS date_vax,
CASE
	WHEN reponse IN ('Etabl. scolaire','Autre','Centre de vaccination de Mamoudzou',
					'Centre de vaccination de Sada','Centre de vaccination de Dembeni','Centre de vaccination de Pamandzi') THEN reponse
	ELSE NULL
END AS lieu_vax,
CASE
	WHEN reponse NOT IN ('Etabl. scolaire','Autre','Centre de vaccination de Mamoudzou',
						'Centre de vaccination de Sada','Centre de vaccination de Dembeni','Centre de vaccination de Pamandzi')
				AND NOT regexp_like(reponse,'^[0-9]{2}/[0-9]{2}/[0-9]{4}$') THEN reponse
	ELSE NULL
END AS vax
from bi_rep_s r
inner join bi_lib_s l on r.nir=l.nir
inner join bi_th_s t on t.nith=l.nith
inner join mouvemen m on m.nisejmouv=t.nisejmouv
inner join sejour s on s.nisejour=m.nisejour
INNER JOIN patient p ON p.nipatient=s.nipatient
inner join ej_srv srv on srv.niservice= m.niservice
where niconcept in (7562,7775,6469,7847,7848,7844)
and dates > '202501000000' AND dates < :datesaisi
and nirsup is not null
AND srv.nom = 'ACTIONS DE SANTE'
order by nda, l.nith, r.nir