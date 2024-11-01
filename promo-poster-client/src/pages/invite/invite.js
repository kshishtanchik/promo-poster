import React, { useState, useEffect } from "react";
import { Container } from "react-bootstrap";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Stack from "react-bootstrap/Stack";
import { useParams } from "react-router-dom";

function Invite() {
  const { activityId } = useParams();
  const [eventData, setEventData] = useState({
    title: "",
    description: "",
    startDate: "",
    activityType: "",
    duration: "",
    address: "",
    periodicity: "",
    baner: null,
    chatId: "",
  });

  useEffect(() => {
    fetch(`http://127.0.0.1:8181/activity/${activityId}`)
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        setEventData(data);
      })
      .catch((err) => {
        console.log(err);
      });
  }, [activityId]);

  return (
    <Container>
      <Card className="text-center">
        <Card.Header>{eventData.activityType}</Card.Header>
       { eventData.baner&&<Card.Img
          variant="top"
          src={eventData.baner}
        />
        }
        <Card.Body>
          <Card.Title>{eventData.title}</Card.Title>
          <Card.Text>
            {eventData.description}
          </Card.Text>
          <Button variant="primary">Записаться</Button>
        </Card.Body>
        <Card.Footer className="text-muted">
          Адресс:
          {eventData.address}
          <Stack direction="horizontal" gap={2}>
            <div className="p-2">Когда: {eventData.startDate}</div>
            <div className="vr  ms-auto" />
            <div className="p-2 ms-auto">Во сколько: {eventData.periodicity}</div>
          </Stack>
        </Card.Footer>
      </Card>
    </Container>
  );
}

export default Invite;
